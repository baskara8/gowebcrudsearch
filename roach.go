package main

import (
"fmt"
"log"
"database/sql"
_ "github.com/lib/pq"
)

func Connect()*sql.DB {
	db, err := sql.Open("postgres", "postgresql://root@localhost:26257/bank?sslmode=disable")

	if err !=nil{
		fmt.Println("Gagal koneksi ke Database `bank`")
		fmt.Print(err.Error())
	}
	return db
}

func main() {
	var db = Connect()
	defer db.Close()
	for i:=2110242;i<=10000000;i++{
	query := fmt.Sprintf("INSERT INTO accounts (id, balance) VALUES (%d, 1000)",i)
	if _, err := db.Exec(query); err != nil {
		log.Fatal(err)
	}
	}
}