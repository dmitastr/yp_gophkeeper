package cryptomanager

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
)

type CryptoManager interface {
	Encrypt([]byte) ([]byte, error)
	Decrypt([]byte) ([]byte, error)
	GenerateKey(int) ([]byte, error)
}

type cryptoManager struct {
	cipher cipher.AEAD
}

func NewCryptoManager(key []byte) (CryptoManager, error) {
	c := &cryptoManager{}
	k := sha256.Sum256(key)
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	c.cipher = gcm

	return c, nil
}

func (c *cryptoManager) Encrypt(token []byte) ([]byte, error) {
	nonce := make([]byte, c.cipher.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := c.cipher.Seal(nil, nonce, token, nil)
	output := append(nonce, ciphertext...)

	return output, nil
}

func (c *cryptoManager) Decrypt(token []byte) ([]byte, error) {
	nonceSize := c.cipher.NonceSize()
	if len(token) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	// Extract nonce and ciphertext
	nonce := token[:nonceSize]
	ciphertext := token[nonceSize:]

	msg, err := c.cipher.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("error decrypting message: %w", err)
	}
	return msg, nil
}

func (c *cryptoManager) GenerateKey(length int) ([]byte, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}
