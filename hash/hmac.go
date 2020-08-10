package hash

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
)

/*
func main() {

	secret := "2dce505d96a53c5768052ee90f3df2055657518dad489160df9913f66042e160"
	data := "33359cee953beb0a2355f15b8d7e2c876107023bcbcbabcfd5c4a262ff7530f2"
	//fmt.Printf("Secret: %s Data: %s\n", secret, data)

	// Get result and encode as hexadecimal string
	sha := hash.New(secret, data)

	//fmt.Println("Result: " + sha)
	fmt.Printf("Result: %s \nLenght: %d\n", sha, len(sha))

	hash2 := "2705728a2b6a2c11a5cdfff15b382ce192bb050ccd65e4fb693595815436a41d"
	ok, _ := hash.Compare(sha, hash2)

	fmt.Printf("Result: %v \n", ok)
}
*/

func New(secret, data string) string {
	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte(secret))

	// Write Data to it
	h.Write([]byte(data))

	// Get result and encode as hexadecimal string
	return hex.EncodeToString(h.Sum(nil))
}

// CompareHexHash to compare passed string as hex string.
func Compare(str, hexStr string) (bool, error) {
	sum, err := hex.DecodeString(str)
	if err != nil {
		return false, err
	}

	orig, err := hex.DecodeString(hexStr)
	if err != nil {
		return false, err
	}

	return subtle.ConstantTimeCompare(sum, orig) == 1, nil
}
