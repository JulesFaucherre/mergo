package git

import (
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBranch(t *testing.T) {
	old, _ := LocalRepository().
		Branch().
		Do(context.Background())
	old = strings.TrimSpace(old)

	branch := "go-test-branch"

	exec.Command("git", "checkout", "-b", branch).Run()

	v, err := LocalRepository().
		Branch().
		Do(context.Background())
	v = strings.TrimSpace(v)

	assert.Equal(t, branch, v)
	assert.Nil(t, err)

	exec.Command("git", "checkout", old).Run()
	exec.Command("git", "branch", "-D", branch).Run()
}

func TestRemote(t *testing.T) {
	key := "test"
	value := "git@gitlab.com:jfaucherre/mergo.git"

	exec.Command("git", "remote", "add", key, value).Run()

	v, err := LocalRepository().
		Remote(key).
		Do(context.Background())
	v = strings.TrimSpace(v)

	assert.Equal(t, v, value)
	assert.Nil(t, err)

	exec.Command("git", "remote", "remove", key).Run()
}
