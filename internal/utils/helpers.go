package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
)

func GetRandomInt(max int64) (int64, error) {
	maxNum := big.NewInt(max + 1)

	randInt, err := rand.Int(rand.Reader, maxNum)
	if err != nil {
		return 0, fmt.Errorf("generating random int error: %w", err)
	}

	return randInt.Int64(), nil
}

func CreateFolder(path string) error {
	if _, err := os.Stat(path); err != nil {
		errMk := os.Mkdir(path, 0750)
		if errMk != nil {
			return fmt.Errorf("failed creating folder with path %s: %w", path, errMk)
		}
	}

	return nil
}

func GetConfigBasePath() string {
	var baseConfigFolder string

	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	baseConfigFolder = fmt.Sprintf("%s/%s", homeDir, ConfigBaseDirName)
	if err := CreateFolder(baseConfigFolder); err != nil {
		panic(err)
	}

	return baseConfigFolder
}

func CheckValInSlice[T comparable](value T, vSlice []T) bool {
	for _, val := range vSlice {
		if val == value {
			return true
		}
	}

	return false
}
