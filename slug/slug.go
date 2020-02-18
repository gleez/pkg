package slug

import (
	"github.com/gosimple/slug"
)

//=============================================================================

// Make returns slug generated from provided string. Will use "en" as language
// substitution.
func New(s string) (sl string) {
	// Maximum of 75 characters for slugs right now
	slug.MaxLength = 75
	return slug.MakeLang(s, "en")
}
