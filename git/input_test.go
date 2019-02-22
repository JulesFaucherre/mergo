package git

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEditor(t *testing.T) {
	v, err := Repository(testRepo).
		GetEditor().
		Do(context.Background())
	v = strings.TrimSpace(v)

	assert.Equal(t, "nano", v)
	assert.Equal(t, nil, err)
}
