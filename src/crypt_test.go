package main

import (
	"testing"
	"bytes"
)

func TestItCanEncryptAndDecryptWithPassphrase(t *testing.T) {
	password := []byte("this is my password")
	data := []byte("this is my data")

	encrypted, err := encrypt(password, data)

	if bytes.Equal(data, encrypted) {
		t.Fatal("Failed to encrypt")
	}

	assertNilError(t, err)

	decrypted, err := decrypt(password, encrypted)

	assertNilError(t, err)

	if ! bytes.Equal(data, decrypted) {
		t.Fatal("Decrypted is not the same as input")
	}
}

func assertNilError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Encountered an error: %v\n", err)
	}
}