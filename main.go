package main

import (
	"flag"
	"fmt"
	"github.com/WadhahJemai/psgen/pkg/utils"
	"strings"
)

func main() {
	length := flag.Int("ln", utils.DefaultLength, "password length")
	upperCase := flag.Bool("u", false, "include upper cases")
	digit := flag.Bool("d", false, "include numbers")
	special := flag.Bool("s", false, "include special characters")
	chars := []string{
		"abcdefghijklmnopqrstuvwxyz",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"~`!@#$%^&*()-_+={}[]|\\;:\"<>,./?",
		"0123456789",
	}

	flag.Parse()

	var strBuilder strings.Builder

	strBuilder.WriteString(chars[0])

	if *upperCase {
		strBuilder.WriteString(chars[1])
	}

	if *digit {
		strBuilder.WriteString(chars[2])
	}

	if *special {
		strBuilder.WriteString(chars[3])
	}

	settings := strBuilder.String()
	var pass strings.Builder

	for i := 0; i < *length; i++ {
		randIndex := utils.GetRandomInt(int64(len(settings)) - 1).Int64()
		pass.WriteRune(rune(settings[randIndex]))
	}

	fmt.Println(pass.String())
}
