package generator

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"github.com/WadhahJemai/psgen/internal/store"
	utils2 "github.com/WadhahJemai/psgen/internal/utils"
	"os"
	"strings"
	"time"
)

func NewCli(store store.Store, encKey string, execTimeout time.Duration) *Cli {
	return &Cli{store: store, key: encKey, execTimeout: execTimeout}
}

func (c *Cli) ExecuteCmd(cmd string, args ...string) string {
	genSet := flag.NewFlagSet("gen", flag.ExitOnError)
	length := genSet.Int("ln", utils2.DefaultLength, "password length")
	upperCase := genSet.Bool("u", false, "include upper cases")
	digits := genSet.Bool("d", false, "include numbers")
	special := genSet.Bool("s", false, "include special characters")
	getSet := flag.NewFlagSet("get", flag.ExitOnError)
	key := getSet.String("key", "", "password key")

	switch cmd {
	case "gen":
		if err := genSet.Parse(args); err != nil {
			return err.Error()
		}
		c.genFlagSet = &GenFlags{
			lowerCase: true,
			upperCase: *upperCase,
			digits:    *digits,
			special:   *special,
			length:    *length,
		}
		err := c.GeneratePassword()
		if err != nil {
			return err.Error()
		}

		return "Password successfully generated and stored"
	case "get":
		if err := getSet.Parse(args); err != nil {
			return err.Error()
		}
		pass, err := c.GetPassword(*key)
		if err != nil {
			return err.Error()
		}
		return pass
	case "export-db":
	case "import-db":
	}

	return "unknown command"
}

func (c *Cli) GeneratePassword() error {
	chars := []string{
		"abcdefghijklmnopqrstuvwxyz",
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"~`!@#$%^&*()-_+={}[]|\\;:\"<>,./?",
		"0123456789",
	}

	toUseIndexes := make([]int, 0)
	toUseIndexes = append(toUseIndexes, 0)

	if c.genFlagSet.upperCase {
		toUseIndexes = append(toUseIndexes, 1)
	}

	if c.genFlagSet.special {
		toUseIndexes = append(toUseIndexes, 2)
	}

	if c.genFlagSet.digits {
		toUseIndexes = append(toUseIndexes, 3)
	}

	var pass strings.Builder

	for i := 0; i < c.genFlagSet.length; i++ {
		index := utils2.GetRandomInt(int64(len(toUseIndexes) - 1))
		targetIndex := toUseIndexes[index]
		charIndex := utils2.GetRandomInt(int64(len(chars[targetIndex]) - 1))

		_, err := pass.WriteString(chars[targetIndex][charIndex : charIndex+1])
		if err != nil {
			return fmt.Errorf("failed to generate password: %w", err)
		}
	}

	generatedPass := pass.String()

	fmt.Println("password: ", generatedPass)
	fmt.Print("Store password? Y[yes] ")

	r := bufio.NewReader(os.Stdin)
	val, _ := r.ReadString('\n')

	val = strings.ToLower(strings.TrimSpace(val))

	if val == "y" || val == "yes" {
		fmt.Print("Give password key: ")

		r := bufio.NewReader(os.Stdin)
		key, _ := r.ReadString('\n')

		encryptedPass, err := utils2.EncryptAES(generatedPass, c.key)
		if err != nil {
			return fmt.Errorf("password not saved. Failed encrypting the generated password: %w", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), c.execTimeout)
		defer cancel()

		if err := c.store.InsertPassword(ctx,
			&store.Password{Key: strings.TrimSpace(key), Value: encryptedPass}); err != nil {
			return fmt.Errorf("failed to store password: %w", err)
		}
	}

	return nil
}

func (c *Cli) GetPassword(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.execTimeout)
	defer cancel()
	p, err := c.store.GetPasswordByKey(ctx, key)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("password with key %s is not present", key)
		} else {
			return "", err
		}
	}

	pass, err := utils2.DecryptAES(p.Value, c.key)
	if err != nil {
		return "", err
	}

	return pass, nil
}
