package generator

import (
	"bufio"
	"context"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"github.com/WadhahJemai/psgen/internal/store"
	"github.com/WadhahJemai/psgen/internal/utils"
	"os"
	"strings"
	"time"
)

func NewCli(store store.Store, cfg *Config) *Cli {
	return &Cli{store: store, config: cfg}
}

func (c *Cli) PrintHelp() {
	fmt.Println(`Usage: psgen <command> -[-]<flags>
Commands:
gen		generates a password
get		retrieves a password from the local sqlite db and prints it out to stdout
export		exports the stored passwords from the local sqlite db to an csv file
import		imports passwords from a csv file into the local sqlite db
help		show help

Use psgen <command> -h or --help for more information about a command.`)

}

func (c *Cli) ExecuteCmd(cmd string, args ...string) string {
	var (
		genFlagSet    *flag.FlagSet
		getFlagSet    *flag.FlagSet
		exportFlagSet *flag.FlagSet
		importFlagSet *flag.FlagSet
	)

	switch cmd {
	case "gen":
		genFlagSet = flag.NewFlagSet("gen", flag.ExitOnError)
		length := genFlagSet.Int("ln", utils.DefaultLength, "password length")
		upperCase := genFlagSet.Bool("u", false, "include upper cases")
		digits := genFlagSet.Bool("d", false, "include numbers")
		special := genFlagSet.Bool("s", false, "include special characters")

		if err := genFlagSet.Parse(args); err != nil {
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
			errW := c.config.WriteLogs([]byte(err.Error()))
			if errW != nil {
				panic(errW)
			}

			return "Error while saving password"
		}

		return "Password successfully generated"
	case "get":
		getFlagSet = flag.NewFlagSet("get", flag.ExitOnError)
		key := getFlagSet.String("key", "", "password key")

		if err := getFlagSet.Parse(args); err != nil {
			return err.Error()
		}

		if len(*key) == 0 {
			getFlagSet.Usage = func() {
				fmt.Println("Missing flags")
				fmt.Println("Usage of get:")
				getFlagSet.PrintDefaults()
			}

			getFlagSet.Usage()

			return ""
		}

		pass, err := c.GetPassword(*key)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return fmt.Sprintf("password with key %s is not present", *key)
			}

			errW := c.config.WriteLogs([]byte(err.Error()))
			if errW != nil {
				panic(errW)
			}

			return fmt.Sprintf("Error while retrieving password with key %s", *key)
		}

		return pass
	case "export":
		exportFlagSet = flag.NewFlagSet("export", flag.ExitOnError)
		outPath := exportFlagSet.String("out", fmt.Sprintf("%s/keys.csv", utils.GetConfigBasePath()), "csv export path")

		if err := exportFlagSet.Parse(args); err != nil {
			return err.Error()
		}

		if err := c.Export(*outPath); err != nil {
			return err.Error()
		}

		return fmt.Sprintf("Passwords exported to %s", *outPath)
	case "import":
		importFlagSet = flag.NewFlagSet("import", flag.ExitOnError)
		csvPath := importFlagSet.String("path", "", "csv path to import (required)")
		concurrentInserts := importFlagSet.Int("c", utils.DefaultConcurrentInserts, "number of concurrent insert operations")
		withEncryption := importFlagSet.Bool("enc", true, "encrypt passwords")

		if err := importFlagSet.Parse(args); err != nil {
			return err.Error()
		}

		if len(*csvPath) == 0 {
			importFlagSet.Usage = func() {
				fmt.Println("Missing flags")
				fmt.Println("Usage of import:")
				importFlagSet.PrintDefaults()
			}

			importFlagSet.Usage()

			return ""
		}

		err := c.Import(*csvPath, *concurrentInserts, *withEncryption)
		if err != nil {
			if errors.Is(err, utils.ErrMalformedCsv) {
				return err.Error()
			}

			errW := c.config.WriteLogs([]byte(err.Error()))
			if errW != nil {
				panic(errW)
			}

			return "Error while importing CSV"
		}

		return "CSV was successfully imported"
	default:
		c.PrintHelp()

		return ""
	}
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
		index, errI := utils.GetRandomInt(int64(len(toUseIndexes) - 1))
		if errI != nil {
			return errI
		}

		targetIndex := toUseIndexes[index]

		charIndex, errCI := utils.GetRandomInt(int64(len(chars[targetIndex]) - 1))
		if errCI != nil {
			return errCI
		}

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

		encryptedPass, err := utils.EncryptAES(generatedPass, c.config.EncKey)
		if err != nil {
			return fmt.Errorf("password not saved. Failed encrypting the generated password: %w", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.config.ExecTimeout)*time.Second)
		defer cancel()

		if err := c.store.InsertPassword(ctx,
			&store.Password{Key: strings.TrimSpace(key), Value: encryptedPass}); err != nil {
			return fmt.Errorf("failed to store password: %w", err)
		}
	}

	return nil
}

func (c *Cli) GetPassword(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.config.ExecTimeout)*time.Second)
	defer cancel()

	p, err := c.store.GetPasswordByKey(ctx, key)
	if err != nil {
		return "", err
	}

	pass, err := utils.DecryptAES(p.Value, c.config.EncKey)
	if err != nil {
		return "", err
	}

	return pass, nil
}

func (c *Cli) Import(csvPath string, concurrency int, withEncryption bool) error {
	content, err := utils.ReadAllCsv(csvPath, ';')
	if err != nil {
		return err
	}

	header := content[0:1]
	body := content[1:]

	if strings.ToLower(header[0][0]) != "key" && strings.ToLower(header[0][1]) != "value" {
		return utils.ErrMalformedCsv
	}

	passwords := make([]*store.Password, 0, len(content[1:]))

	if withEncryption {
		for _, row := range body {
			encPass, err := utils.EncryptAES(row[1], c.config.EncKey)
			if err != nil {
				return err
			}

			passwords = append(passwords, &store.Password{Key: row[0], Value: encPass})
		}
	} else {
		for _, row := range body {
			passwords = append(passwords, &store.Password{Key: row[0], Value: row[1]})
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.config.ExecTimeout)*time.Second)
	defer cancel()

	if err := c.store.BatchInsertPassword(ctx, passwords, concurrency); err != nil {
		return err
	}

	return nil
}

func (c *Cli) Export(path string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(c.config.ExecTimeout)*time.Second)
	defer cancel()

	p, err := c.store.GetAllPasswords(ctx)
	if err != nil {
		return err
	}

	header := []string{"Key", "Value"}
	body := make([][]string, 0, len(p))

	for _, pass := range p {
		decPass, err := utils.DecryptAES(pass.Value, c.config.EncKey)
		if err != nil {
			return err
		}
		body = append(body, []string{pass.Key, decPass})
	}

	if err := utils.WriteToCsv(header, body, path, ';'); err != nil {
		return err
	}

	return nil
}
