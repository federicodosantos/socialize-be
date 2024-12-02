package mysql

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func DBInit() *sqlx.DB {
	username := "root"
	password := ""
	host := "db"
	port := "3306"
	database := "db-socialize"

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, host, port, database)
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}

	return db
}
