package functest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGET(t *testing.T) {

	resp := httpGET(t)
	assert.Equal(t, 200, resp.StatusCode())

}

func TestPOST(t *testing.T) {

	resp := httpPOST(t)
	assert.Equal(t, 200, resp.StatusCode())

}
