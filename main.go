package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/Passengers")

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Database opened")

	defer db.Close()
}
