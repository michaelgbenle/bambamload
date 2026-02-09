package utils

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"errors"
)

// Encrypt encrypts the provided text using AES encryption with the specified key and returns the base64-encoded result.
// The encryptionKey must be 16, 24, or 32 bytes in length. Returns an error if the key is invalid or if encryption fails.
func Encrypt(encryptionKey, text string) (string, error) {
	key := []byte(encryptionKey)
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", errors.New("invalid key length; must be 16, 24, or 32 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	bs := block.BlockSize() // 16 for AES
	data := []byte(text)

	// Apply PKCS#7 padding correctly
	// If length is already multiple of block size, add a full block of padding
	padding := bs - (len(data) % bs)
	if padding == bs { // This happens when len(data)%bs == 0
		padding = bs
	}

	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	data = append(data, padText...)

	// Encrypt using ECB mode (manual loop since Go doesn't provide ECB directly)
	ciphertext := make([]byte, len(data))
	for i := 0; i < len(data); i += bs {
		block.Encrypt(ciphertext[i:i+bs], data[i:i+bs])
	}

	// Return as URL-safe base64 string
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a base64-encoded ciphertext string using the provided AES encryption key in ECB mode.
// The encryption key must be 16, 24, or 32 bytes long to match AES requirements.
// Returns the decrypted plaintext as a string or an error if decryption fails.
func Decrypt(encryptionKey, ciphertextStr string) (string, error) {
	if len(ciphertextStr) == 0 || ciphertextStr == "" {
		return "", errors.New("ciphertext is empty")
	}
	// Decode the base64 string
	ciphertext, err := base64.URLEncoding.DecodeString(ciphertextStr)
	if err != nil {
		return "", err
	}

	// Validate key length
	key := []byte(encryptionKey)
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return "", errors.New("invalid key length; must be 16, 24, or 32 bytes")
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Validate ciphertext length (must be a multiple of the block size)
	bs := block.BlockSize()
	if len(ciphertext)%bs != 0 {
		return "", errors.New("ciphertext length must be a multiple of block size")
	}

	// Decrypt using ECB mode
	plaintext := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i += bs {
		block.Decrypt(plaintext[i:i+bs], ciphertext[i:i+bs])
	}

	// Remove padding (PKCS#5/PKCS#7)
	padding := int(plaintext[len(plaintext)-1])
	if padding > bs || padding <= 0 {
		return "", errors.New("invalid padding")
	}
	return string(plaintext[:len(plaintext)-padding]), nil
}
