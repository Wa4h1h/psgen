package main

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/WadhahJemai/psgen/internal/store"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

func main() {
	d := store.NewDatabase("psgen.db", 5*time.Second)

	if err := d.InitSchema(); err != nil {
		panic(err.Error())
	}

	err := d.BatchInsertPassword(context.Background(), []*store.Password{
		{Key: "1", Value: "1"},
		{Key: "2", Value: "2"},
		{Key: "3", Value: "3"},
	}, 100)
	if err != nil {
		fmt.Println(err.Error())
	}
}
