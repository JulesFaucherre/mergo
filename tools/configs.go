// The package tools provides utilities for mergo
// The main utility you will find is to save, get and delete user's credentials
// You should use these functions when managing user's credentials since they
// come with their way of handling encryption
package tools

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
	configDir = path.Join(os.Getenv("HOME"), ".config", "mergo")
	keyPath   = path.Join(configDir, ".key")
	key       []byte
)

// GetHostConfig returns the config stored for the host 'host'
func GetHostConfig(host string) ([]byte, error) {
	p := path.Join(configDir, host)
	s, err := os.Stat(p)

	if err != nil {
		return nil, err
	}
	if s.IsDir() {
		return nil, fmt.Errorf("%s must not be a directory", p)
	}

	content, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}

	key, err := getEncryptionKey()
	if err != nil {
		return nil, err
	}

	return decrypt(content, key)
}

// WriteHostConfig writes the config 'content' for the host 'host'
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

	key, err := getEncryptionKey()
	if err != nil {
		return err
	}

	encrypted, err := encrypt(content, key)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(p, encrypted, 0644)
}

// DeleteHostConfig deletes the config for host 'host'
func DeleteHostConfig(host string) error {
	return os.Remove(path.Join(configDir, host))
}

func encrypt(plaintext []byte, key []byte) ([]byte, error) {
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

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
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

func getEncryptionKey() ([]byte, error) {
	if key == nil {
		_, err := os.Stat(keyPath)
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		if os.IsNotExist(err) {
			return createEncryptionKey()
		}
		return ioutil.ReadFile(keyPath)
	}
	return key, nil
}

func createEncryptionKey() ([]byte, error) {
	key = make([]byte, keySize, keySize)
	_, err := rand.Read(key)

	err = ioutil.WriteFile(keyPath, key, 0400)
	if err != nil {
		return nil, err
	}

	return key, nil
}
