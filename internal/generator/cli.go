package generator

import (
	"flag"
	"github.com/WadhahJemai/psgen/internal/store"
	utils2 "github.com/WadhahJemai/psgen/internal/utils"
	"strings"
)

type Flags struct {
	lowerCase bool
	upperCase bool
	digits    bool
	special   bool
	length    int
}

type Cli struct {
	flags *Flags
	store store.Store
}

func NewCli() *Cli {
	return &Cli{}
}

func (c *Cli) ParseFlags() {
	length := flag.Int("ln", utils2.DefaultLength, "password length")
	upperCase := flag.Bool("u", false, "include upper cases")
	digits := flag.Bool("d", false, "include numbers")
	special := flag.Bool("s", false, "include special characters")

	flag.Parse()

	c.flags = &Flags{
		lowerCase: true,
		upperCase: *upperCase,
		digits:    *digits,
		special:   *special,
		length:    *length,
	}
}

func (c *Cli) GeneratePassword() string {
	chars := []string{
		"abcdefghijklmnopqrstuvwxyz",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"~`!@#$%^&*()-_+={}[]|\\;:\"<>,./?",
		"0123456789",
	}
	targetBound := 0

	if c.flags.upperCase {
		targetBound++
	}

	if c.flags.digits {
		targetBound++
	}

	if c.flags.special {
		targetBound++
	}

	var pass strings.Builder

	for i := 0; i < c.flags.length; i++ {
		char := utils2.GetRandomInt(int64(targetBound))
		charRandIndex := utils2.GetRandomInt(int64(len(chars[char])) - 1)
		pass.WriteString(chars[char][charRandIndex : charRandIndex+1])
	}

	return pass.String()
}
