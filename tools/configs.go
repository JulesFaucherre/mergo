package tools

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

var (
	configDir = path.Join(os.Getenv("HOME"), ".config", "mergo")
)

func GetHostConfig(host string) ([]byte, error) {
	p := path.Join(configDir, host)
	s, err := os.Stat(p)

	if err != nil {
		return []byte{}, err
	}
	if s.IsDir() {
		return []byte{}, fmt.Errorf("%s must not be a directory", p)
	}

	return ioutil.ReadFile(p)
}

func WriteHostConfig(host string, content []byte) error {
	s, err := os.Stat(configDir)

	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if !os.IsNotExist(err) && !s.IsDir() {
		return fmt.Errorf("%s must be a directory", configDir)
	}
	if os.IsNotExist(err) {
		if err = os.MkdirAll(configDir, 0755); err != nil {
			return err
		}
	}

	p := path.Join(configDir, host)

	return ioutil.WriteFile(p, content, 0644)
}
