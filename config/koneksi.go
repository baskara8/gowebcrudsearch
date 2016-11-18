package config

import "fmt"
import "database/sql"
// import _"github.com/denisenkom/go-mssqldb" import mssql
import _ "github.com/go-sql-driver/mysql" 


func Connect()*sql.DB {
	dbDriver := "mysql"
	dbUser   := "root" 
	// dbURL    := "localhost" 
	dbPass   := ""
	dbName   := "latihan"
	// db, err := sql.Open(dbDriver,"server="+dbURL+";user id="+dbUser+";password="+dbPass+";database="+dbName+"") Open SQL Server
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)

	if err !=nil{
		fmt.Println("Gagal koneksi ke Database `dbpksm`")
		fmt.Print(err.Error())
	}
	return db
}