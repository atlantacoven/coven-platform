package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	a := NewTest(t)

	res, err := a.Request("GET", "/", nil)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "test", res.Data["environment"])
}
