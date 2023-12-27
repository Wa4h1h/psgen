package main

import (
	_ "embed"
	"fmt"
	"github.com/WadhahJemai/psgen/internal/generator"
	"github.com/WadhahJemai/psgen/internal/store"
	"github.com/WadhahJemai/psgen/internal/utils"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"time"
)

var (
	DbPath        = utils.GetEnv[string](utils.DbPathVar, "psgen.db", false)
	ExecTimeout   = utils.GetEnv[uint](utils.ExecTimeout, "5", false)
	EncryptionKey = utils.GetEnv[string](utils.EncryptKeyVar, "", true)
)

func main() {
	d := store.NewDatabase(DbPath, "EncryptionKey", time.Duration(ExecTimeout)*time.Second)
	defer func() {
		if err := d.Close(); err != nil {
			panic(err.Error())
		}
	}()

	if err := d.InitSchema(); err != nil {
		panic(err.Error())
	}

	c := generator.NewCli(d.Q, EncryptionKey, time.Duration(ExecTimeout)*time.Second)
	args := os.Args

	if len(args) < 3 {
		fmt.Println("expected 'gen' or 'get' or export-db or import-db subcommands")
		os.Exit(1)
	}

	result := c.ExecuteCmd(args[1], args[2:]...)

	fmt.Println(result)
}
