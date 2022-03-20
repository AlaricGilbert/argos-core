package bitcoin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLookupBTCNetwork(t *testing.T) {
	ips, err := LookupBTCNetwork()
	assert.Nil(t, err)
	if len(ips) > 0 {
		t.Logf("Get bitcoin seeds ok:[%v]", ips)
	}
}
