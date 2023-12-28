package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
)

func GenerateNewRandomKey() (string, error) {
	buf := make([]byte, 32)

	_, err := rand.Read(buf)
	if err != nil {
		return "", fmt.Errorf("creating new key failed: %w", err)
	}

	return base64.StdEncoding.EncodeToString(buf), nil
}

func EncryptAES(str string, key string) (string, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", fmt.Errorf("decoding base64 key error: %w", err)
	}

	aesC, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("creating new cipher error: %w", err)
	}

	gcm, err := cipher.NewGCM(aesC)
	if err != nil {
		return "", fmt.Errorf("creating gcm  error: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())

	_, err = rand.Read(nonce)
	if err != nil {
		return "", fmt.Errorf("reading random bytes for nonce error: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(str), nil)

	return hex.EncodeToString(ciphertext), nil
}

func DecryptAES(str string, key string) (string, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", fmt.Errorf("decoding base64 key error: %w", err)
	}

	aesC, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("creating new cipher error: %w", err)
	}

	gcm, err := cipher.NewGCM(aesC)
	if err != nil {
		return "", fmt.Errorf("creating gcm  error: %w", err)
	}

	decStr, err := hex.DecodeString(str)
	if err != nil {
		return "", fmt.Errorf("decoding hey cipher text error: %w", err)
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := decStr[:nonceSize], decStr[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("deciphering screte error: %w", err)
	}

	return string(plaintext), nil
}
