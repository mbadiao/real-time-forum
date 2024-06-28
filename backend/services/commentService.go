package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"realTimeForum/backend/controllers"
	"realTimeForum/backend/models"
	"sync"

	"github.com/gorilla/websocket"
)

var comment models.Comment
var processedMessages = make(map[string]bool)

var connections = make(map[*websocket.Conn]bool)
var connectionsLock sync.Mutex

func CommentWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()

	connectionsLock.Lock()
	connections[conn] = true
	connectionsLock.Unlock()

	defer func() {
		connectionsLock.Lock()
		delete(connections, conn)
		connectionsLock.Unlock()
	}()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}
		messageKey := string(message)
		// if processedMessages[messageKey] {
		// 	continue
		// }

		err = json.Unmarshal(message, &comment)
		if err != nil {
			fmt.Println("Error unmarshaling data:", err)
			return
		}
		fmt.Println("comment", comment)
		posts, errs := controllers.HandleCommentData(w, r, comment)
		ServicePosts.Posts = posts
		ServicePosts.Errors = errs

		// newCommentResponse := map[string]interface{}{
		// 	"type":   "new_comment",
		// 	"Posts":  posts,
		// 	"Errors": errs,
		// }

		username, _, _ := controllers.GetUsernameAndFirstname(comment.UserID)

		newComment := map[string]interface{}{
			"comment": comment,
			"user":    username,
			"type":    "new_comment",
		}

		// response, err := json.Marshal(newCommentResponse)
		responses, err := json.Marshal(newComment)
		fmt.Println("newCommentResponse", newComment)
		if err != nil {
			fmt.Println("Error marshaling response:", err)
			return
		}

		broadcastMessaget(messageType, responses)

		processedMessages[messageKey] = true
	}
}

func broadcastMessaget(messageType int, message []byte) {
	connectionsLock.Lock()
	defer connectionsLock.Unlock()
	for conn := range connections {
		if err := conn.WriteMessage(messageType, message); err != nil {
			fmt.Println("Error writing message:", err)
			conn.Close()
			delete(connections, conn)
		}
	}
}
