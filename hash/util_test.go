package hash

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/pbkdf2"
)

const password = "abcdefghijkl"
const salt = "505d96a53c5768052ee90f3df"

func TestGenerateSHA1(t *testing.T) {
	generateHexTest(t, sha1.New)
}

func TestGenerateSHA256(t *testing.T) {
	generateHexTest(t, sha256.New)
	generateHash(t)
}

func TestGenerateSHA512(t *testing.T) {
	generateHexTest(t, sha512.New)
}

func TestCompareSHA1(t *testing.T) {
	compareHexTest(t, sha1.New)
}

func TestCompareSHA256(t *testing.T) {
	compareHexTest(t, sha256.New)
	compareHash(t)
}

func TestCompareSHA512(t *testing.T) {
	compareHexTest(t, sha512.New)
}

func generateHexTest(t *testing.T, hash func() hash.Hash) {
	hexStr, salt, err := NewHex(hash, password)
	if err != nil {
		t.Fatalf("Get error on HexHash(): %s", err.Error())
	}

	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		t.Fatalf("Get error on hex.DecodeString(hex): %s", err.Error())
	}

	saltBytes, err := hex.DecodeString(salt)
	if err != nil {
		t.Fatalf("Get error on hex.DecodeString(salt): %s", err.Error())
	}

	sum := pbkdf2.Key([]byte(password), saltBytes, Iter, KeyLength, hash)

	Equal(t, sum, decoded, "got %s, expected %s", string(decoded), password)
}

func compareHexTest(t *testing.T, hash func() hash.Hash) {
	hex, salt, err := NewHex(hash, password)
	if err != nil {
		t.Fatalf("Get error on HexHash(): %s", err.Error())
	}
	got, err := CompareHex(hash, password, hex, salt)
	if err != nil {
		t.Fatalf("Get error on CompareHexHash(): %s", err.Error())
	}

	Equal(t, true, got, "got %t, expected %t", got, true)
}

func generateHash(t *testing.T) {
	want := "77cbc57dc6ce8bbfb63d5883d86d3a55f7437cdcc6f604490436df34603f5aec"
	hexStr := New(salt, password)

	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		t.Fatalf("Get error on hex.DecodeString(hex): %s", err.Error())
	}

	Equal(t, want, hexStr, "got %s, expected %s", string(decoded), password)
}

func compareHash(t *testing.T) {
	hash1 := "77cbc57dc6ce8bbfb63d5883d86d3a55f7437cdcc6f604490436df34603f5aec"
	hash2 := "77cbc57dc6ce8bbfb63d5883d86d3a55f7437cdcc6f604490436df34603f5aec"

	ok, err := Compare(hash1, hash2)
	if err != nil {
		t.Fatalf("Get error on compare: %s", err.Error())
	}

	Equal(t, true, ok, "got %v, expected %v", ok, true)
}

func Equal(t *testing.T, expected interface{}, actual interface{}, format string, args ...interface{}) {
	assert.Equal(t, expected, actual, fmt.Sprintf(format, args...))
}
