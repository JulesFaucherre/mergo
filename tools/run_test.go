package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasic(t *testing.T) {
	c, err := Run([][]string{
		{"ls"},
		{"grep", "run"},
		{"wc", "-l"},
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, "       2\n", c)
}
