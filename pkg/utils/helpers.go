package utils

import (
	"crypto/rand"
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
