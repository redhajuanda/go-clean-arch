package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCeiling(t *testing.T) {
	got := Ceiling(1499, 500)
	assert.Equal(t, 1500, got)
}
