package vault

import (
	"bytes"
	"crypto/rand"
	"errors"
	"testing"
)

func kek(t *testing.T, id string) *StaticKEK {
	t.Helper()
	k := make([]byte, 32)
	rand.Read(k)
	p, err := NewStaticKEK(id, k)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestSealOpenRoundTrip(t *testing.T) {
	kp := kek(t, "kek-1")
	secret := []byte("S3cr3t-root-password!")
	s, err := Seal(kp, "acc_root@db-prod-01", 1, secret)
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(s.Ciphertext, secret) || bytes.Contains(s.WrappedDEK, secret) {
		t.Fatal("plaintext visible in stored data") // T-V1
	}
	got, err := Open(kp, s)
	if err != nil || !bytes.Equal(got, secret) {
		t.Fatalf("round trip failed: %v", err)
	}
}

func TestUniqueDEKPerSecret(t *testing.T) {
	kp := kek(t, "kek-1")
	a, _ := Seal(kp, "acc", 1, []byte("same"))
	b, _ := Seal(kp, "acc", 1, []byte("same"))
	if bytes.Equal(a.Ciphertext, b.Ciphertext) || bytes.Equal(a.WrappedDEK, b.WrappedDEK) {
		t.Fatal("ciphertexts must differ (fresh DEK and nonce)")
	}
}

func TestWrongKEKFails(t *testing.T) {
	s, _ := Seal(kek(t, "kek-1"), "acc", 1, []byte("x"))
	if _, err := Open(kek(t, "kek-1"), s); !errors.Is(err, ErrDecrypt) {
		t.Fatalf("want ErrDecrypt, got %v", err)
	}
}

func TestSwappedRowFails(t *testing.T) {
	kp := kek(t, "kek-1")
	s, _ := Seal(kp, "acc_A", 1, []byte("password-of-A"))
	s.AccountID = "acc_B" // attacker copies A's row to B
	if _, err := Open(kp, s); !errors.Is(err, ErrDecrypt) {
		t.Fatalf("AAD must bind ciphertext to account, got %v", err)
	}
}

func TestTamperedCiphertextFails(t *testing.T) {
	kp := kek(t, "kek-1")
	s, _ := Seal(kp, "acc", 1, []byte("secret"))
	s.Ciphertext[len(s.Ciphertext)-1] ^= 0x01
	if _, err := Open(kp, s); !errors.Is(err, ErrDecrypt) {
		t.Fatalf("want ErrDecrypt, got %v", err)
	}
}

func TestRewrapKEKRotation(t *testing.T) { // T-V6
	oldKP, newKP := kek(t, "kek-1"), kek(t, "kek-2")
	s, _ := Seal(oldKP, "acc", 3, []byte("rotate-me"))
	r, err := Rewrap(oldKP, newKP, s)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(r.Ciphertext, s.Ciphertext) {
		t.Fatal("ciphertext should be unchanged by rewrap")
	}
	got, err := Open(newKP, r)
	if err != nil || string(got) != "rotate-me" {
		t.Fatalf("open after rewrap: %v", err)
	}
	if _, err := Open(oldKP, r); err == nil {
		t.Fatal("old KEK must not open rewrapped secret")
	}
}

func TestKEKLength(t *testing.T) {
	if _, err := NewStaticKEK("x", make([]byte, 16)); err == nil {
		t.Fatal("expected error for 128-bit KEK")
	}
}
