package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
)

func GetRandomInt(max int64) int64 {
	maxNum := big.NewInt(max + 1)

	randInt, err := rand.Int(rand.Reader, maxNum)
	if err != nil {
		panic(err)
	}

	return randInt.Int64()
}

func GetConfigBasePath() string {
	var baseConfigFolder string

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	baseConfigFolder = fmt.Sprintf("%s/%s", homeDir, ConfigBaseDirName)
	if _, err := os.Stat(baseConfigFolder); err != nil {
		errMk := os.Mkdir(baseConfigFolder, 0750)
		if errMk != nil {
			panic(errMk)
		}
	}

	return baseConfigFolder
}

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
		return "", err
	}

	aesC, err := aes.NewCipher(keyBytes)
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
	keyBytes, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return "", err
	}

	aesC, err := aes.NewCipher(keyBytes)
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
