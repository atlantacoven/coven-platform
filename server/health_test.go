package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHealthCheck(t *testing.T) {
	s := TestServer(t)

	res, err := s.Request("GET", "/", nil)
	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "test", res.Data["environment"])
}
