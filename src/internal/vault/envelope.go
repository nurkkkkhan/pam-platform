// Package vault implements envelope encryption for privileged credentials
// (ADR-003, FR-V1, FR-V6). Only Go standard-library primitives are used.
package vault

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

// KeyProvider wraps and unwraps data-encryption keys with a key-encryption
// key (KEK) held outside the database (file/env now; HSM/KMS later).
type KeyProvider interface {
	KeyID() string
	Wrap(dek, aad []byte) ([]byte, error)
	Unwrap(wrapped, aad []byte) ([]byte, error)
}

// StaticKEK is a KeyProvider backed by a 32-byte key loaded from a file or
// environment variable. Suitable for the test bench only.
type StaticKEK struct {
	id  string
	gcm cipher.AEAD
}

// NewStaticKEK creates a provider from a 32-byte key.
func NewStaticKEK(id string, key []byte) (*StaticKEK, error) {
	if len(key) != 32 {
		return nil, errors.New("KEK must be 32 bytes (AES-256)")
	}
	g, err := newGCM(key)
	if err != nil {
		return nil, err
	}
	return &StaticKEK{id: id, gcm: g}, nil
}

func (k *StaticKEK) KeyID() string                        { return k.id }
func (k *StaticKEK) Wrap(dek, aad []byte) ([]byte, error) { return seal(k.gcm, dek, aad) }
func (k *StaticKEK) Unwrap(w, aad []byte) ([]byte, error) { return open(k.gcm, w, aad) }

// Sealed is what is stored in the `secrets` table. No plaintext, no DEK.
type Sealed struct {
	AccountID  string
	Version    int
	KEKID      string
	WrappedDEK []byte
	Ciphertext []byte
}

// ErrDecrypt is returned for any authentication failure (wrong key, swapped
// row, modified ciphertext). Details are not leaked.
var ErrDecrypt = errors.New("vault: decryption failed")

func aad(accountID string, version int) []byte {
	return []byte(fmt.Sprintf("pam-platform/v1|%s|%d", accountID, version))
}

// Seal encrypts a secret with a fresh DEK and wraps the DEK with the KEK.
func Seal(kp KeyProvider, accountID string, version int, secret []byte) (*Sealed, error) {
	dek := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, dek); err != nil {
		return nil, err
	}
	defer zero(dek)
	g, err := newGCM(dek)
	if err != nil {
		return nil, err
	}
	a := aad(accountID, version)
	ct, err := seal(g, secret, a)
	if err != nil {
		return nil, err
	}
	w, err := kp.Wrap(dek, a)
	if err != nil {
		return nil, err
	}
	return &Sealed{AccountID: accountID, Version: version, KEKID: kp.KeyID(), WrappedDEK: w, Ciphertext: ct}, nil
}

// Open decrypts a sealed secret. Caller must zero the result after use.
func Open(kp KeyProvider, s *Sealed) ([]byte, error) {
	if s.KEKID != kp.KeyID() {
		return nil, fmt.Errorf("vault: secret wrapped with KEK %q, provider has %q", s.KEKID, kp.KeyID())
	}
	a := aad(s.AccountID, s.Version)
	dek, err := kp.Unwrap(s.WrappedDEK, a)
	if err != nil {
		return nil, ErrDecrypt
	}
	defer zero(dek)
	g, err := newGCM(dek)
	if err != nil {
		return nil, err
	}
	pt, err := open(g, s.Ciphertext, a)
	if err != nil {
		return nil, ErrDecrypt
	}
	return pt, nil
}

// Rewrap moves a secret to a new KEK without touching the ciphertext (FR-V6).
func Rewrap(oldKP, newKP KeyProvider, s *Sealed) (*Sealed, error) {
	a := aad(s.AccountID, s.Version)
	dek, err := oldKP.Unwrap(s.WrappedDEK, a)
	if err != nil {
		return nil, ErrDecrypt
	}
	defer zero(dek)
	w, err := newKP.Wrap(dek, a)
	if err != nil {
		return nil, err
	}
	out := *s
	out.KEKID = newKP.KeyID()
	out.WrappedDEK = w
	return &out, nil
}

func newGCM(key []byte) (cipher.AEAD, error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(b)
}

// seal returns nonce || ciphertext||tag.
func seal(g cipher.AEAD, pt, aad []byte) ([]byte, error) {
	nonce := make([]byte, g.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return g.Seal(nonce, nonce, pt, aad), nil
}

func open(g cipher.AEAD, in, aad []byte) ([]byte, error) {
	if len(in) < g.NonceSize() {
		return nil, ErrDecrypt
	}
	return g.Open(nil, in[:g.NonceSize()], in[g.NonceSize():], aad)
}

func zero(b []byte) {
	for i := range b {
		b[i] = 0
	}
}
