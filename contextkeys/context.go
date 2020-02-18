// Package contextkeys stores the keys to the context accessor
// functions, letting generated code safely set values in contexts
// without exposing the setters to the outside world.
package contextkeys

// KeyString should be used when setting and fetching context values
type KeyString string

// Added useful keys
const (
	Tenant = KeyString("X-Tenant-ID")
	User   = KeyString("X-User-ID")

	TokenPointer = KeyString("tokenPointer")
)
