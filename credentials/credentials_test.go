package credentials

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHostConfig(t *testing.T) {
	host := "test_host"
	content := []byte("test content")

	err := WriteHostConfig(host, content)
	assert.Nil(t, err)

	c, _ := ioutil.ReadFile(path.Join(configDir, host))
	assert.NotEqual(t, c, content)

	c, err = GetHostConfig(host)
	assert.Nil(t, err)
	assert.Equal(t, c, content)

	err = DeleteHostConfig(host)
	assert.Nil(t, err)
}
