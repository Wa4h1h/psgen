package main

import (
	_ "embed"
	"fmt"
	"os"
	"time"

	"github.com/Wa4h1h/psgen/internal/generator"
	"github.com/Wa4h1h/psgen/internal/store"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	cfg, err := generator.LoadConfig()
	if err != nil {
		panic(fmt.Errorf("failed loading config: %w", err))
	}

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

	if len(args) < 2 {
		c.PrintHelp()
		os.Exit(0)
	}

	result := c.ExecuteCmd(args[1], args[2:]...)

	fmt.Print(result)
}
