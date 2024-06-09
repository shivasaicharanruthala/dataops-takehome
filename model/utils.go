package model

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// encrypt encrypts plaintext using AES encryption with the provided key.
func encrypt(plaintext string) (*string, error) {
	key := []byte("example key 1234")

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	plaintextBytes := []byte(plaintext)

	// Padding plaintext if necessary
	padding := aes.BlockSize - len(plaintextBytes)%aes.BlockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	plaintextBytes = append(plaintextBytes, padText...)

	ciphertext := make([]byte, len(plaintextBytes))

	mode := cipher.NewCBCEncrypter(block, make([]byte, aes.BlockSize))
	mode.CryptBlocks(ciphertext, plaintextBytes)

	encryptedText := base64.StdEncoding.EncodeToString(ciphertext)

	return &encryptedText, nil
}

// decrypt decrypts ciphertext using AES decryption with the provided key.
func decrypt(ciphertextStr string) (*string, error) {
	key := []byte("example key 1234")

	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextStr)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	decrypted := make([]byte, len(ciphertext))

	mode := cipher.NewCBCDecrypter(block, make([]byte, aes.BlockSize))
	mode.CryptBlocks(decrypted, ciphertext)

	// Remove padding
	padding := int(decrypted[len(decrypted)-1])
	decrypted = decrypted[:len(decrypted)-padding]

	plainText := string(decrypted)

	return &plainText, nil
}
