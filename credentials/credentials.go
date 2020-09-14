// The package credentials provides encryption for mergo
// The main utility you will find is to save, get and delete user's credentials
// You should use these functions when managing user's credentials since they
// come with their way of handling encryption
package credentials

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

const keySize = 32

var (
	configDir string
	key       []byte

	ErrNoHostConfig = errors.New("no config for given host")
)

// the init function ensures the config directory is present and creates it if
// it does not exist
func init() {
	config, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}

	// ensure the configuration directory exists
	configDir = path.Join(config, "mergo")
	s, err := os.Stat(configDir)

	if os.IsNotExist(err) {
		err = os.MkdirAll(configDir, 0755)
	}
	if err != nil {
		panic(err)
	}
	if s != nil && !s.IsDir() {
		panic(fmt.Errorf("%s must be a directory", configDir))
	}

	// create or load key
	keyPath := path.Join(configDir, ".key")
	_, err = os.Stat(keyPath)
	if os.IsNotExist(err) {
		key, err = createEncryptionKey(keyPath)
	} else if err == nil {
		key, err = ioutil.ReadFile(keyPath)
	}
	if err != nil {
		panic(err)
	}
}

// GetHostConfig returns the credentials stored for the host 'host'
func GetHostConfig(host string) ([]byte, error) {
	p := path.Join(configDir, host)
	_, err := os.Stat(p)

	if os.IsNotExist(err) {
		return nil, ErrNoHostConfig
	}
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	return decrypt(content)
}

// WriteHostConfig writes the config 'content' as the credentials for the host
// 'host'
func WriteHostConfig(host string, content []byte) error {
	p := path.Join(configDir, host)

	encrypted, err := encrypt(content)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(p, encrypted, 0644)
}

// DeletedHostConfig takes a host name and deletes its credentials
func DeleteHostConfig(host string) error {
	err := os.Remove(path.Join(configDir, host))
	if os.IsNotExist(err) {
		return ErrNoHostConfig
	}
	return err
}

func encrypt(plaintext []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(ciphertext []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func createEncryptionKey(keyPath string) ([]byte, error) {
	key = make([]byte, keySize, keySize)
	_, err := rand.Read(key)

	err = ioutil.WriteFile(keyPath, key, 0400)
	if err != nil {
		return nil, err
	}

	return key, nil
}
