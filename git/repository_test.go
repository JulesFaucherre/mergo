package git

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	pwd, _   = os.Getwd()
	testRepo = path.Join(pwd, "..", "tests")
)

func TestBranch(t *testing.T) {
	v, err := Repository(testRepo).
		Branch().
		Do(context.Background())

	assert.Equal(t, v, "test\n")
	assert.Equal(t, err, nil)
}

func TestRemote(t *testing.T) {
	v, err := Repository(testRepo).
		Remote("origin").
		Do(context.Background())

	assert.Equal(t, v, "git@gitlab.com:jfaucherre/mergo.git\n")
	assert.Equal(t, err, nil)
}
