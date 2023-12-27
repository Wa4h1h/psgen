package main

import (
	_ "embed"
	"fmt"
	"github.com/WadhahJemai/psgen/internal/generator"
	"github.com/WadhahJemai/psgen/internal/store"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"time"
)

func main() {
	cfg := generator.LoadConfig()

	d := store.NewDatabase(cfg.DbPath, time.Duration(cfg.ExecTimeout)*time.Second)
	defer func() {
		if err := d.Close(); err != nil {
			panic(err.Error())
		}
	}()

	if err := d.InitSchema(); err != nil {
		panic(err.Error())
	}

	c := generator.NewCli(d.Q, cfg)
	args := os.Args

	if len(args) < 3 {
		fmt.Println("expected 'gen' or 'get' or export-db or import-db subcommands")
		os.Exit(1)
	}

	result := c.ExecuteCmd(args[1], args[2:]...)

	fmt.Println(result)
}
