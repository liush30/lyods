package db

import (
	"database/sql"
	"log"
)

func GetDb() *sql.DB {
	db, err := sql.Open("mysql", "root:lyods@123@tcp(192.168.1.212:3306)/sit_nf_vaw")
	if err != nil {
		log.Fatal(err)
	}
	return db
}
