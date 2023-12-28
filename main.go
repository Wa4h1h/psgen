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

	if !utils.CheckValInSlice(args[1], []string{"gen", "get", "import", "export"}) {
		fmt.Println("expected 'gen' or 'get' or 'export' or 'import' commands")
		os.Exit(1)
	}

	if len(args) < 3 {
		fmt.Println("flags missing")
		os.Exit(1)
	}

	result := c.ExecuteCmd(args[1], args[2:]...)

	fmt.Println(result)
}
