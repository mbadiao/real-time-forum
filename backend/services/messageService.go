package services

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"realTimeForum/backend/controllers"
	"realTimeForum/backend/utils"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var connMess = make(map[int]*websocket.Conn)

func WsMessages(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	receiver := query.Get("receiver")
	// fmt.Println("query: ", query)
	// fmt.Println("receiver: ", receiver)
	// fmt.Println("connected: ", connected)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("probleme lors de l'initialisation: ", err)
		return
	}
	cookie := controllers.GetCookieHandler(w, r)
	connectedUser, err := utils.GetUserIDFromDB(db, cookie)
	connMess[connectedUser] = conn
	if err != nil {
		fmt.Println("error retrieving user ID from the database: ", err)
		return
	}
	// conn.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(connectedUser)))
	receiverInt, errAtoi := strconv.Atoi(receiver)
	// fmt.Println("receiver int: ", receiverInt)
	if errAtoi != nil {
		fmt.Println(errAtoi)
		return
	}
	conversations, errConv := FetchConversation(db, connectedUser, receiverInt)
	if errConv != nil {
		fmt.Println(errConv)
		return
	}
	conn.WriteMessage(websocket.TextMessage, conversations)
	go Reader(conn, connectedUser, receiverInt)
}

func Reader(conn *websocket.Conn, sender, receiver int) {
	defer conn.Close()
	for {
		_, message, err := conn.ReadMessage()
		// fmt.Println("message", string(message))
		if err != nil {
			fmt.Println("error reading received message: ", err)
			return
		}

		if len(message) != 0 {
			messageChannel <- SenderReceiver{Sender: sender, Receiver: receiver}
		}

		conversations_id, errCOnv := storeConversation(db, sender, receiver)
		
		if errCOnv != nil {
			fmt.Println(errCOnv)
			return
		}

		u1,_:=GetUsernameByID(db,sender)
		u2,_:=GetUsernameByID(db,receiver)

		if controllers.Interactions[u1]== nil{
			controllers.Interactions[u1] = make(map[string]time.Time)
		}
		controllers.Interactions[u1][u2]= time.Now()

		if controllers.Interactions[u2]== nil{
			controllers.Interactions[u2] = make(map[string]time.Time)
		}
		controllers.Interactions[u2][u1]= time.Now()

		controllers.Broadcast <- "ok"

		errDb := storeMessage(db, conversations_id, sender, receiver, string(message))
		if errDb != nil {
			fmt.Println(errDb)
			return
		}
		conversations, errConv := FetchConversation(db, sender, receiver)
		if errConv != nil {
			fmt.Println(errConv)
			return
		}
		conn.WriteMessage(websocket.TextMessage, conversations)
		receiverConn, ok := connMess[receiver]
		if !ok {
			fmt.Println("Receiver not connected")
			continue
		}

		if receiverConn != nil {
			receiverConn.WriteMessage(websocket.TextMessage, conversations)
		}
	}
}

func storeConversation(db *sql.DB, sender, receiver int) (int64, error) {

	var conversationID int64
	err := db.QueryRow("SELECT conversation_id FROM Conversations WHERE (participant1_id = ? AND participant2_id = ?) OR (participant1_id = ? AND participant2_id = ?)", sender, receiver, receiver, sender).Scan(&conversationID)
	if err != nil {
		if err == sql.ErrNoRows {
			result, err := db.Exec("INSERT INTO Conversations (participant1_id, participant2_id) VALUES (?, ?)", sender, receiver)
			if err != nil {
				return 0, err
			}
			conversationID, err = result.LastInsertId()
			if err != nil {
				return 0, err
			}
		} else {
			return 0, err
		}
	}

	return conversationID, nil
}

func storeMessage(db *sql.DB, conversationID int64, sender int, receiver int, content string) error {
	_, err := db.Exec("INSERT INTO Messages (conversation_id, sender_id, receiver_id, content) VALUES (?, ?, ?, ?)", conversationID, sender, receiver, content)
	if err != nil {
		return err
	}

	return nil
}

func FetchConversation(db *sql.DB, sender, receiver int) ([]byte, error) {
	rows, err := db.Query("SELECT sender_id, receiver_id, content, DATE(creation_date) FROM Messages WHERE conversation_id IN (SELECT conversation_id FROM Conversations WHERE (participant1_id = ? AND participant2_id = ?) OR (participant1_id = ? AND participant2_id = ?)) ORDER BY creation_date ASC ", sender, receiver, receiver, sender)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	conversation := make([]map[string]string, 0)
	for rows.Next() {
		var senderID, receiverID, content, date string
		err := rows.Scan(&senderID, &receiverID, &content, &date)
		if err != nil {
			return nil, err
		}
		formatSender, formatReceiver, errUser := GetUsernames(db, senderID, receiverID)
		if errUser != nil {
			return nil, errUser
		}
		formatSender += " " + date
		formatReceiver += " " + date

		message := map[string]string{
			"sender":         senderID,
			"receiver":       receiverID,
			"content":        content,
			"formatSender":   "by " + formatSender,
			"formatReceiver": "by " + formatReceiver,
		}
		conversation = append(conversation, message)
	}

	jsonData, err := json.Marshal(conversation)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func GetUsernames(db *sql.DB, sender, receiver string) (string, string, error) {
	var username1, username2 string

	senderInt, err := strconv.Atoi(sender)
	receiverInt, err := strconv.Atoi(sender)

	err = db.QueryRow("SELECT username FROM Users WHERE user_id = ?", senderInt).Scan(&username1)
	if err != nil {
		return "", "", err
	}

	err = db.QueryRow("SELECT username FROM Users WHERE user_id = ?", receiverInt).Scan(&username2)
	if err != nil {
		return "", "", err
	}

	return username1, username2, nil
}

func GetUsernameByID(db *sql.DB, userID int) (string, error) {
    var username string
    query := "SELECT username FROM Users WHERE user_id = ?"
     
    // Exécuter la requête
    err := db.QueryRow(query, userID).Scan(&username)
    if err != nil {
        return "", err
    }
    return username, nil
}