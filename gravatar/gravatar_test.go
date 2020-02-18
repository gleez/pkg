package gravatar

import (
	"fmt"
	"log"
	"testing"
)

func TestAvatar(t *testing.T) {
	expected := "https://www.gravatar.com/avatar/d9290cc27176c6fc74f4002f40fc9db8?s=256"
	actual := Avatar("philipp@zoonman.com", 256)

	if actual != expected {
		log.Print(fmt.Errorf("Actual value \"%s\" doesn't match expected \"%s\"", actual, expected))
		t.Fail()
	}

}
