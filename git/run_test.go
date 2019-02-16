package git

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimple(t *testing.T) {
	c, err := run(context.Background(), [][]string{
		{"git", "config", "--global", "core.editor", "nvim"},
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, "", c)
}

func TestMany(t *testing.T) {
	c, err := run(context.Background(), [][]string{
		{"ls"},
		{"grep", "run"},
		{"wc", "-l"},
	})

	assert.Equal(t, nil, err)
	assert.Equal(t, "       2\n", c)
}
