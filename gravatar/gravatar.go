package gravatar

import (
	"crypto/md5"
	"fmt"
)

// Create Gravatar picture URL from email address
func Avatar(email string, size uint) string {
	hash := md5.Sum([]byte(email))
	return fmt.Sprintf("https://www.gravatar.com/avatar/%x?s=%d", hash, size)
}
