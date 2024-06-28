package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"realTimeForum/backend/controllers"
	"realTimeForum/backend/utils"
	"strconv"

	"github.com/gorilla/websocket"
)

type TypingMessage struct {
	ID     string `json:"id"`
	Typing string `json:"typing`
}

var typingConnect = make(map[int]*websocket.Conn)

func TypingHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("probleme lors de l'initialisation: ", err)
		return
	}

	cookie := controllers.GetCookieHandler(w, r)
	connectedUser, err := utils.GetUserIDFromDB(db, cookie)
	typingConnect[connectedUser] = conn
	if err != nil {
		fmt.Println("erreur lors de l'ajout: ", err)
		return
	}
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Erreur lors de la lecture du message: ", err)
				return
			}

			var msg TypingMessage
			err = json.Unmarshal(message, &msg)
			if err != nil {
				fmt.Println("Erreur lors de la désérialisation du message: ", err)
				return
			}

			name, errName := fetchUserFirstname(connectedUser)
			if errName != nil {
				return
			}
			receiverint, err00 := strconv.Atoi(msg.ID)
			if err00 != nil {
				return
			}
			receiverConn, ok := typingConnect[receiverint]
			if !ok {
				fmt.Println("Receiver not connected")
				continue
			}
			if receiverConn != nil {
				receiverConn.WriteMessage(websocket.TextMessage, []byte(name))
			}
			if msg.Typing == "false" {
				receiverConn.WriteMessage(websocket.TextMessage, []byte{})
			}
		}
	}()
}

func fetchUserFirstname(userID int) (string, error) {
	var firstname, lastname string
	query := "SELECT firstname, lastname FROM Users WHERE user_id = ?"
	err := db.QueryRow(query, userID).Scan(&firstname, &lastname)
	if err != nil {
		return "", err
	}
	return firstname + " " + lastname, nil
}
