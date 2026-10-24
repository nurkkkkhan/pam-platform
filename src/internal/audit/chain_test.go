package audit

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func sampleChain(t *testing.T, n int) *Chain {
	t.Helper()
	c := NewChain()
	base := time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC)
	for i := 0; i < n; i++ {
		if _, err := c.Append(base.Add(time.Duration(i)*time.Second), "alice", "session.command", "session/S-1",
			map[string]any{"cmd": "ls -la", "i": i}); err != nil {
			t.Fatal(err)
		}
	}
	return c
}

func TestVerifyIntactChain(t *testing.T) {
	c := sampleChain(t, 50)
	if err := Verify(c.Events(), c.Head()); err != nil {
		t.Fatalf("intact chain should verify: %v", err)
	}
}

func TestDetectModification(t *testing.T) {
	ev := sampleChain(t, 20).Events()
	ev[9].Details = json.RawMessage(`{"cmd":"echo harmless","i":9}`)
	assertBrokenAt(t, Verify(ev, ""), 10)
}

func TestDetectActorChange(t *testing.T) {
	ev := sampleChain(t, 5).Events()
	ev[2].Actor = "bob"
	assertBrokenAt(t, Verify(ev, ""), 3)
}

func TestDetectDeletion(t *testing.T) {
	ev := sampleChain(t, 20).Events()
	ev = append(ev[:4], ev[5:]...)
	assertBrokenAt(t, Verify(ev, ""), 5)
}

func TestDetectReorder(t *testing.T) {
	ev := sampleChain(t, 10).Events()
	ev[3], ev[4] = ev[4], ev[3]
	assertBrokenAt(t, Verify(ev, ""), 4)
}

func TestDetectRehashedForgery(t *testing.T) {
	// Attacker modifies event 3 and recomputes its hash: the next link breaks.
	ev := sampleChain(t, 6).Events()
	ev[2].Action = "session.noop"
	h, _ := ComputeHash(ev[2])
	ev[2].Hash = h
	assertBrokenAt(t, Verify(ev, ""), 4)
}

func TestDetectTailTruncationWithAnchor(t *testing.T) {
	c := sampleChain(t, 10)
	ev := c.Events()[:8]
	if err := Verify(ev, ""); err != nil {
		t.Fatalf("truncated prefix is internally consistent: %v", err)
	}
	if err := Verify(ev, c.Head()); !errors.Is(err, ErrAnchorMismatch) {
		t.Fatalf("want ErrAnchorMismatch, got %v", err)
	}
}

func TestDetailsWhitespaceIrrelevant(t *testing.T) {
	ev := sampleChain(t, 3).Events()
	ev[1].Details = json.RawMessage(`{ "i" : 1 , "cmd" : "ls -la" }`)
	if err := Verify(ev, ""); err != nil {
		t.Fatalf("canonicalization should ignore formatting: %v", err)
	}
}

func TestDetectLargeIntegerModification(t *testing.T) {
	c := NewChain()
	_, err := c.Append(time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC), "alice", "session.query", "session/S-1", map[string]any{
		"rows_affected": json.Number("9007199254740993"),
	})
	if err != nil {
		t.Fatal(err)
	}
	events := c.Events()
	if err := Verify(events, c.Head()); err != nil {
		t.Fatalf("original event should verify: %v", err)
	}
	events[0].Details = json.RawMessage(`{"rows_affected":9007199254740992}`)
	assertBrokenAt(t, Verify(events, c.Head()), 1)
}

func TestReturnedEventsCannotMutateChain(t *testing.T) {
	c := NewChain()
	returned, err := c.Append(time.Date(2026, 10, 24, 9, 0, 0, 0, time.UTC), "alice", "session.command", "session/S-1", map[string]any{"cmd": "ls"})
	if err != nil {
		t.Fatal(err)
	}
	returned.Details[0] = '['
	if err := Verify(c.Events(), c.Head()); err != nil {
		t.Fatalf("mutating Append result changed stored event: %v", err)
	}

	exported := c.Events()
	exported[0].Details[0] = '['
	if err := Verify(c.Events(), c.Head()); err != nil {
		t.Fatalf("mutating Events result changed stored event: %v", err)
	}
}

func assertBrokenAt(t *testing.T, err error, seq uint64) {
	t.Helper()
	var ve *VerifyError
	if !errors.As(err, &ve) {
		t.Fatalf("expected VerifyError, got %v", err)
	}
	if ve.Seq != seq {
		t.Fatalf("expected break at seq %d, got %d (%s)", seq, ve.Seq, ve.Reason)
	}
}
