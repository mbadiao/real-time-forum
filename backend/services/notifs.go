package services

import (
	"fmt"
	"net/http"
	"realTimeForum/backend/controllers"
	"realTimeForum/backend/utils"

	"github.com/gorilla/websocket"
)

// Définir la structure et le canal
type SenderReceiver struct {
	Sender   int
	Receiver int
}

var messageChannel = make(chan SenderReceiver)

var notifConn = make(map[int]*websocket.Conn)

func Notifshandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("probleme lors de l'initialisation: ", err)
		return
	}
	cookie := controllers.GetCookieHandler(w, r)
	connectedUser, err := utils.GetUserIDFromDB(db, cookie)
	fmt.Println("connected", connectedUser)
	notifConn[connectedUser] = conn
	if err != nil {
		fmt.Println("error retrieving user ID from the database: ", err)
		return
	}
	for {
		select {
		case sr := <-messageChannel:
			var SenderName string
			err = db.QueryRow("SELECT username FROM Users WHERE user_id = ?", sr.Sender).Scan(&SenderName)
			if err != nil {
				fmt.Println("error retrieving name from the database: ", err)
				return
			}
			receiverConn, ok := notifConn[sr.Receiver]
			if !ok {
				fmt.Println("user not connected")
				continue
			}

			if receiverConn != nil {
				err = receiverConn.WriteMessage(websocket.TextMessage, []byte(SenderName))
				if err != nil {
					delete(notifConn, sr.Receiver)
					fmt.Println("error sending message: ", err)
				}
			}
		}
	}
}
