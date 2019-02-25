package git

import (
	"context"
	"strings"
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
	c = strings.TrimSpace(c)

	assert.Equal(t, nil, err)
	assert.Equal(t, "2", c)
}
