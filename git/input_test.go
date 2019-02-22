package git

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEditor(t *testing.T) {
	exec.Command("git", "config", "core.editor", "nano").Run()

	v, err := LocalRepository().
		GetEditor().
		Do(context.Background())
	v = strings.TrimSpace(v)

	assert.Equal(t, "nano", v)
	assert.Nil(t, nil, err)

	exec.Command("git", "config", "--unset", "core.editor").Run()
}
