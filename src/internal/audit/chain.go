// Package audit implements the tamper-evident (hash-chained) audit log
// described in ADR-005 and requirement FR-AU2.
//
//	hash_n = SHA-256(hash_{n-1} || canonical(event_n))
//
// Any modification, deletion or reordering of events breaks the chain.
package audit

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// GenesisHash is the prev_hash of the first event (32 zero bytes, hex).
var GenesisHash = hex.EncodeToString(make([]byte, sha256.Size))

// Event is one audit record.
type Event struct {
	Seq      uint64          `json:"seq"`
	TS       time.Time       `json:"ts"`
	Actor    string          `json:"actor"`
	Action   string          `json:"action"`
	Object   string          `json:"object"`
	Details  json.RawMessage `json:"details,omitempty"`
	PrevHash string          `json:"prev_hash"`
	Hash     string          `json:"hash"`
}

// payload is the canonical content that is hashed (everything except Hash).
type payload struct {
	Seq      uint64          `json:"seq"`
	TS       string          `json:"ts"`
	Actor    string          `json:"actor"`
	Action   string          `json:"action"`
	Object   string          `json:"object"`
	Details  json.RawMessage `json:"details,omitempty"`
	PrevHash string          `json:"prev_hash"`
}

func canonical(e Event) ([]byte, error) {
	var details json.RawMessage
	if len(e.Details) > 0 {
		// Re-marshal details so that whitespace differences do not matter.
		var v any
		if !json.Valid(e.Details) {
			return nil, errors.New("details: invalid JSON")
		}
		dec := json.NewDecoder(bytes.NewReader(e.Details))
		dec.UseNumber()
		if err := dec.Decode(&v); err != nil {
			return nil, fmt.Errorf("details: %w", err)
		}
		b, err := json.Marshal(v) // map keys are sorted by encoding/json
		if err != nil {
			return nil, err
		}
		details = b
	}
	return json.Marshal(payload{
		Seq: e.Seq, TS: e.TS.UTC().Format(time.RFC3339Nano), Actor: e.Actor,
		Action: e.Action, Object: e.Object, Details: details, PrevHash: e.PrevHash,
	})
}

// ComputeHash returns the chain hash for e (using e.PrevHash).
func ComputeHash(e Event) (string, error) {
	c, err := canonical(e)
	if err != nil {
		return "", err
	}
	prev, err := hex.DecodeString(e.PrevHash)
	if err != nil {
		return "", fmt.Errorf("prev_hash: %w", err)
	}
	h := sha256.New()
	h.Write(prev)
	h.Write(c)
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Chain stores audit events in memory.
type Chain struct {
	events []Event
	head   string
}

// NewChain returns an empty chain.
func NewChain() *Chain { return &Chain{head: GenesisHash} }

// Append links a new event to the chain and returns the stored copy.
func (c *Chain) Append(ts time.Time, actor, action, object string, details any) (Event, error) {
	var raw json.RawMessage
	if details != nil {
		b, err := json.Marshal(details)
		if err != nil {
			return Event{}, err
		}
		raw = b
	}
	e := Event{Seq: uint64(len(c.events)) + 1, TS: ts, Actor: actor, Action: action,
		Object: object, Details: raw, PrevHash: c.head}
	h, err := ComputeHash(e)
	if err != nil {
		return Event{}, err
	}
	e.Hash = h
	c.events = append(c.events, e)
	c.head = h
	return cloneEvent(e), nil
}

// Events returns a copy of all events.
func (c *Chain) Events() []Event {
	if len(c.events) == 0 {
		return nil
	}
	result := make([]Event, len(c.events))
	for i, e := range c.events {
		result[i] = cloneEvent(e)
	}
	return result
}

func cloneEvent(e Event) Event {
	e.Details = append(json.RawMessage(nil), e.Details...)
	return e
}

// Head returns the current head hash (to be anchored externally).
func (c *Chain) Head() string { return c.head }

// VerifyError reports the first broken position.
type VerifyError struct {
	Seq    uint64
	Reason string
}

func (e *VerifyError) Error() string {
	return fmt.Sprintf("audit chain broken at seq %d: %s", e.Seq, e.Reason)
}

// ErrAnchorMismatch is returned when the chain is internally valid but its
// head differs from the externally anchored hash (e.g. tail truncation).
var ErrAnchorMismatch = errors.New("head hash does not match anchor")

// Verify checks continuity of seq numbers, prev_hash links and hashes.
// If anchor is non-empty, the final head must equal it.
func Verify(events []Event, anchor string) error {
	prev := GenesisHash
	for i, e := range events {
		want := uint64(i) + 1
		if e.Seq != want {
			return &VerifyError{Seq: want, Reason: fmt.Sprintf("expected seq %d, found %d (deleted or reordered event)", want, e.Seq)}
		}
		if e.PrevHash != prev {
			return &VerifyError{Seq: e.Seq, Reason: "prev_hash does not link to previous event"}
		}
		h, err := ComputeHash(e)
		if err != nil {
			return &VerifyError{Seq: e.Seq, Reason: err.Error()}
		}
		if h != e.Hash {
			return &VerifyError{Seq: e.Seq, Reason: "content hash mismatch (event modified)"}
		}
		prev = e.Hash
	}
	if anchor != "" && anchor != prev {
		return ErrAnchorMismatch
	}
	return nil
}
