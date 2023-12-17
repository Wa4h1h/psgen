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
	defer func() {
		if err := d.Close(); err != nil {
			fmt.Println(err.Error())
		}
	}()

	if err := d.BatchInsertPassword(context.Background(), []*store.Password{
		{Key: "1", Value: "2"},
		{Key: "2", Value: "2"},
		{Key: "3", Value: "3"},
		{Key: "4", Value: "4"},
		{Key: "5", Value: "5"},
		{Key: "6", Value: "6"},
		{Key: "7", Value: "7"},
		{Key: "8", Value: "8"},
		{Key: "9", Value: "9"},
		{Key: "10", Value: "10"},
		{Key: "11", Value: "11"},
		{Key: "12", Value: "12"},
		{Key: "13", Value: "13"},
	}); err != nil {
		panic(err.Error())
	}
}
