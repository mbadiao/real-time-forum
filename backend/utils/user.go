package utils

import (
	"database/sql"
	"fmt"
	"net/http"
)

func ScanWithSessions(db *sql.DB, query string, params ...interface{}) ([]*Session, error) {
	rows, err := db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sessions []*Session
	for rows.Next() {
		var session Session
		if err := rows.Scan(&session.SessionID, &session.UserID, &session.Cookie_value, &session.ExpirationDate); err != nil {
			return nil, err
		}
		sessions = append(sessions, &session)
	}
	return sessions, nil
}

// func GetCurrentUserId(w http.ResponseWriter, r *http.Request, db *sql.DB) (int, error) {
// 	var userID int
// 	cookie := GetCookieHandler(w, r)
// 	query := `SELECT user_id FROM Sessions WHERE cookie_value = ?`
// 	err := db.QueryRow(query, cookie).Scan(&userID)
// 	if err != nil {
// 		return 0, fmt.Errorf("failed to get user ID from session: %v", err)
// 	}
// 	return userID, nil
// }

func GetCookie(r *http.Request, db *sql.DB) string {
	cookie, err := r.Cookie("ForumCookie")
	if err != nil {
		return ""
	}
	return (cookie.Value)
}

func GetUserIDFromDB(db *sql.DB, cookie string) (int, error) {
	fmt.Println("COOKIE", cookie)
	stmt, err := db.Prepare("SELECT user_id FROM Sessions WHERE cookie_value = ?")
	if err != nil {
		fmt.Println("CHOIX 1")
		return 0, err
	}
	defer stmt.Close()

	var userID int
	err = stmt.QueryRow(cookie).Scan(&userID)
	if err != nil {
		fmt.Println("CHOIX 2")
		return 0, err
	}

	return userID, nil
}

