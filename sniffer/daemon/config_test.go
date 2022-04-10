package daemon

import (
	"testing"
)

func TestRandIdentifier(t *testing.T) {
	t.Log(randIdentifier())
}
