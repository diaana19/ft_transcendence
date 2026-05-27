package utils

import (
	"testing"
)

func TestHashString_ProducesHash(t *testing.T) {
	hash, err := HashString("mypassword")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hash == "" {
		t.Fatal("hash should not be empty")
	}
	if hash == "mypassword" {
		t.Fatal("hash should not equal the original string")
	}
}

func TestHashString_DifferentHashesForSameInput(t *testing.T) {
	hash1, _ := HashString("mypassword")
	hash2, _ := HashString("mypassword")
	if hash1 == hash2 {
		t.Fatal("bcrypt should produce different hashes for the same input (random salt)")
	}
}

func TestCheckHashString_CorrectPassword(t *testing.T) {
	hash, _ := HashString("mypassword")
	if !CheckHashString("mypassword", hash) {
		t.Fatal("correct password should match its hash")
	}
}

func TestCheckHashString_WrongPassword(t *testing.T) {
	hash, _ := HashString("mypassword")
	if CheckHashString("wrongpassword", hash) {
		t.Fatal("wrong password should not match hash")
	}
}

func TestCheckHashString_EmptyPassword(t *testing.T) {
	hash, _ := HashString("mypassword")
	if CheckHashString("", hash) {
		t.Fatal("empty password should not match hash")
	}
}

func TestHashString_EmptyInput(t *testing.T) {
	hash, err := HashString("")
	if err != nil {
		t.Fatalf("hashing empty string should not error: %v", err)
	}
	if !CheckHashString("", hash) {
		t.Fatal("empty string should match its own hash")
	}
}
