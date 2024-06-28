package utils

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mattn/go-sqlite3"
)

type Session struct {
	SessionID      int
	UserID         int
	Cookie_value   string
	ExpirationDate time.Time
}

type Generic struct {
	Type string
	Data json.RawMessage
}

func AnalyzeError(err error) string {
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == sqlite3.ErrConstraint {
			if strings.Contains(sqliteErr.Error(), "UNIQUE constraint failed: Users.username") {
				fmt.Println("erreur : le nom d'utilisateur existe déjà")
				return "username"
			} else if strings.Contains(sqliteErr.Error(), "UNIQUE constraint failed: Users.email") {
				fmt.Println("erreur : l'adresse email existe déjà")
				return "email"
			}
		}
		fmt.Printf("erreur de contrainte unique : %v\n", sqliteErr)
		return "other1"
	}
	fmt.Printf("erreur SQLite : %v\n", err)
	return "other2"
}
