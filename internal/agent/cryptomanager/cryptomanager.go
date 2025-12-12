package cryptomanager

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
)

type CryptoManager interface {
	Encrypt(*string, []byte) (string, error)
	Decrypt(*string, []byte) (string, error)
}

type cryptoManager struct {
	cipher cipher.AEAD
}

func NewCryptoManager() CryptoManager {
	c := &cryptoManager{}
	return c
}

func (c *cryptoManager) createCipher(key []byte) error {
	k := sha256.Sum256(key)
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	c.cipher = gcm
	return nil
}

func (c *cryptoManager) Encrypt(token *string, key []byte) (string, error) {
	if c.cipher == nil {
		if err := c.createCipher(key); err != nil {
			return "", err
		}
	}

	nonce := make([]byte, c.cipher.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := c.cipher.Seal(nil, nonce, []byte(*token), nil)
	output := append(nonce, ciphertext...)

	msg := base64.StdEncoding.EncodeToString(output)
	fmt.Printf("Writing token=%s\n", msg)

	return msg, nil
}

func (c *cryptoManager) Decrypt(token *string, key []byte) (string, error) {
	fmt.Printf("Decrypting token=%s\n", *token)

	if c.cipher == nil {
		if err := c.createCipher(key); err != nil {
			return "", err
		}
	}
	data, err := base64.StdEncoding.DecodeString(*token)

	nonceSize := c.cipher.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	// Extract nonce and ciphertext
	nonce := data[:nonceSize]
	ciphertext := data[nonceSize:]

	msg, err := c.cipher.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("error decrypting message: %w", err)
	}
	return string(msg), nil
}
