package controllers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"realTimeForum/data"
	"sort"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func GetUsernameAndId(cookie string) (int, string, error) {
	db := data.CreateTable()

	var userID int
	query := `SELECT user_id FROM Sessions WHERE cookie_value = ?`
	err := db.QueryRow(query, cookie).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", fmt.Errorf("no session found for the given cookie value")
		}
		return 0, "", fmt.Errorf("failed to get user ID from session: %v", err)
	}

	var username string
	query = `SELECT username FROM Users WHERE user_id = ?`
	err = db.QueryRow(query, userID).Scan(&username)
	if err != nil {
		if err == sql.ErrNoRows {
			return userID, "", fmt.Errorf("no user found for the given user ID")
		}
		return userID, "", fmt.Errorf("failed to get username from user ID: %v", err)
	}

	return userID, username, nil
}

func UpdateUserStatus(username string, status bool) error {
	db := data.CreateTable()
	query := `UPDATE Users SET user_status = ? WHERE username = ?`
	result, err := db.Exec(query, status, username)
	if err != nil {
		return fmt.Errorf("failed to update user status: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to retrieve affected rows: %v", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with the given ID")
	}

	return nil
}

func FetchUsers() error {
	db := data.CreateTable()

	query := `SELECT username, user_status, firstname, lastname, user_id FROM Users`
	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("failed to query users: %v", err)
	}
	defer rows.Close()

	users = make(map[string]User)

	for rows.Next() {
		var username, firstname, lastname string
		var status bool
		var userID int
		if err := rows.Scan(&username, &status, &firstname, &lastname, &userID); err != nil {
			return fmt.Errorf("failed to scan row: %v", err)
		}
		users[username] = User{
			Username:  username,
			Online:    status,
			Firstname: firstname,
			Lastname:  lastname,
			UserID:    userID,
		}
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("rows iteration error: %v", err)
	}

	return nil
}

type User struct {
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Online    bool   `json:"online"`
	UserID    int    `json:"user_id"`
}

var (
	users       = make(map[string]User)
	usersMutex  = sync.Mutex{}
	connections = make(map[*websocket.Conn]string)
	Broadcast    = make(chan string)
	Interactions = make(map[string]map[string]time.Time) // Les interactions entre utilisateurs
)

var StatusUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleOnlineUser(w http.ResponseWriter, r *http.Request) {
	ws, err := StatusUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	var cookie string
	err = ws.ReadJSON(&cookie)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	// fmt.Println("cookie", cookie)

	_, username, _ := GetUsernameAndId(cookie)

	// fmt.Println("username", username)

	usersMutex.Lock()
	connections[ws] = username
	UpdateUserStatus(username, true)
	usersMutex.Unlock()

	Broadcast <- "ok"

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}

	usersMutex.Lock()
	delete(connections, ws)
	UpdateUserStatus(username, false)
	usersMutex.Unlock()

	Broadcast <- "ok"
}

func HandleMessages() {
	for {
		// fmt.Println("Dans la boucle handle messs")
		// Get the updated user status from the broadcast channel
		user := <-Broadcast
		FetchUsers()
		fmt.Println(user)
		usersMutex.Lock()
		// Notify all clients about the updated list of users
		for conn := range connections {
			// fmt.Println("users map", users)

			currentUsername := connections[conn]
			filteredUsers := filterUsers(currentUsername)

			err := conn.WriteJSON(filteredUsers)
			if err != nil {
				conn.Close()
				delete(connections, conn)
			}
		}
		usersMutex.Unlock()
	}
}

func filterUsers(currentUsername string) []User {
	var userList []User
	for username, user := range users {
		if username != currentUsername {
			userList = append(userList, user)
		}
	}

	sort.SliceStable(userList, func(i, j int) bool {
		interactionTimeI, existsI := Interactions[currentUsername][userList[i].Username]
		interactionTimeJ, existsJ := Interactions[currentUsername][userList[j].Username]

		if existsI && existsJ {
			return interactionTimeI.After(interactionTimeJ)
		}
		if existsI {
			return true
		}
		if existsJ {
			return false
		}

		return userList[i].Username < userList[j].Username
	})

	return userList
}
