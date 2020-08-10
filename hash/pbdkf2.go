package hash

/*
https://github.com/Code-Hex/saltissimo

func main() {
    gotFromForm := "password"
    // 1. Code to generate hash
    hash, key, err := saltissimo.HexHash(sha256.New, gotFromForm)
    if err != nil {
        panic(err)
    }
    // *Code to save some values

    // 2. Code to compare hash
    // *Code to retrieve the value from a database etc.
    // *Assume that it has already been substituted.
    isSame, err := saltissimo.CompareHexHash(sha256.New, gotFromForm, hash, key)
    if err != nil {
        panic(err)
    }
    if isSame {
        fmt.Println("Hello user!!")
    } else {
        fmt.Println("Who are you...?")
    }
}

*/

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"hash"

	"golang.org/x/crypto/pbkdf2"
)

// SaltLength specifies the length of a random byte sequence.
// KeyLength, Iter for pbkdf2.Key arguments.
var (
	SaltLength = 36
	KeyLength  = 512
	Iter       = 10000
)

// HexHash to generate PBDKF2 as Hex string.
// returns PBDKF2, secret key, error
func NewHex(hash func() hash.Hash, str string) (string, string, error) {
	key, err := RandomBytes(SaltLength)
	if err != nil {
		return "", "", err
	}
	return PBDKF2Hex(hash, str, key), hex.EncodeToString(key), nil
}

// PBDKF2Hex creates a hex string from PBDKF2 as its name
func PBDKF2Hex(hash func() hash.Hash, str string, key []byte) string {
	b := pbkdf2.Key([]byte(str), key, Iter, KeyLength, hash)
	return hex.EncodeToString(b)
}

// CompareHexHash to compare passed string and PBDKF2 as hex string.
func CompareHex(hash func() hash.Hash, str, hexStr, key string) (bool, error) {
	kb, err := hex.DecodeString(key)
	if err != nil {
		return false, err
	}

	orig, err := hex.DecodeString(hexStr)
	if err != nil {
		return false, err
	}

	sum := pbkdf2.Key([]byte(str), kb, Iter, KeyLength, hash)
	return subtle.ConstantTimeCompare(sum, orig) == 1, nil
}

// RandomBytes generate a random byte slice.
func RandomBytes(l int) ([]byte, error) {
	b := make([]byte, l)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}
