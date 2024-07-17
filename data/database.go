package data

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func CreateTable() *sql.DB {
	_, errNofile := os.Stat("./data/forum.db")

	db, err := sql.Open("sqlite3", "./data/forum.db")
	if err != nil {
		log.Println(err.Error())
	}
	if errNofile != nil {
		sqlcode, err := os.ReadFile("./data/table.sql")
		if err != nil {
			log.Println(err.Error())
		}
		_, err = db.Exec(string(sqlcode))

		if err != nil {
			log.Println(err.Error())
		}
	}
	return db
}

func GeneratePrepare(text string) string {
	nb := len(strings.Split(text, ","))
	a := strings.Repeat("?,", nb)
	return "(" + a[:len(a)-1] + ")"
}

func Insert(db *sql.DB, table, values string, data ...interface{}) {
	prepare := GeneratePrepare(values)
	Query := fmt.Sprintf("INSERT INTO %v %v values %v", table, values, prepare)
	insert, err := db.Prepare(Query)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	_, err = insert.Exec(data...)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
