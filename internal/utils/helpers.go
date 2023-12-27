package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"math/big"
)

func GetRandomInt(max int64) int64 {
	maxNum := big.NewInt(max + 1)

	randInt, err := rand.Int(rand.Reader, maxNum)
	if err != nil {
		panic(err)
	}

	return randInt.Int64()
}

func EncryptAES(str string, key string) (string, error) {
	aesC, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(aesC)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(str), nil)

	return hex.EncodeToString(ciphertext), nil

}

func DecryptAES(str string, key string) (string, error) {
	aesC, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err.Error())
	}

	gcm, err := cipher.NewGCM(aesC)
	if err != nil {
		panic(err.Error())
	}

	decStr, err := hex.DecodeString(str)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := decStr[:nonceSize], decStr[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		panic(err)
	}

	return string(plaintext), nil
}
