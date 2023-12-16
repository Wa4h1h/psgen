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

	err := d.InsertPassword(context.Background(), &store.Password{Key: "ccc", Value: "ddd"})
	if err != nil {
		fmt.Println(err.Error())
	}

	p, err := d.GetPasswordByKey(context.Background(), "aaa")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(p)

}
