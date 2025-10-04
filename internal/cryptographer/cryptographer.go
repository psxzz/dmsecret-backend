package cryptographer

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

type encrypter struct {
	key []byte
}

func New(key string) (*encrypter, error) {
	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return nil, fmt.Errorf("error decoding key: %w", err)
	}

	if !isKeyValid(keyBytes) {
		return nil, errors.New("invalid crypto key")
	}

	return &encrypter{key: keyBytes}, nil
}

func (e *encrypter) Encrypt(payload string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("couldn't create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("couldn't create AES GCM cipher: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("couldn't create nonce: %w", err)
	}

	encrypted := gcm.Seal(nonce, nonce, []byte(payload), nil)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func (e *encrypter) Decrypt(encoded string) (string, error) {
	encrypted, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("couldn't decode base64: %w", err)
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("couldn't create AES cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("couldn't create AES GCM cipher: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(encrypted) < nonceSize {
		return "", errors.New("encrypted payload is too short")
	}

	payload, err := gcm.Open(nil, encrypted[:nonceSize], encrypted[nonceSize:], nil)
	if err != nil {
		return "", fmt.Errorf("couldn't decrypt payload: %w", err)
	}

	return string(payload), nil
}

func isKeyValid(key []byte) bool {
	switch len(key) {
	case 16, 24, 32:
		return true
	default:
	}

	return false
}
