package hasher

import (
	"crypto/sha512"
	"fmt"
	"os"
	"testing"
)

func hashWithSalt(str, salt string) string {
	hash := sha512.New()
	hash.Write([]byte(str))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
func TestHash(t *testing.T) {
	hasher := NewHasher(os.Getenv("Salt"))

	// 1
	str := "hhh"
	expectedHash := hashWithSalt(str, os.Getenv("Salt"))
	hashedStr, err := hasher.Hash(str)

	if err != nil {
		t.Errorf("Error hash: %v", err)
	}

	if hashedStr != expectedHash {
		t.Errorf("Excepted %s, result %s", expectedHash, hashedStr)
	}

	//2
	str = ""
	expectedHash = hashWithSalt(str, os.Getenv("Salt"))
	hashedStr, err = hasher.Hash(str)

	if err != nil {
		t.Errorf("Error hash: %v", err)
	}

	if hashedStr != expectedHash {
		t.Errorf("Excepted %s, result %s", expectedHash, hashedStr)
	}
}
