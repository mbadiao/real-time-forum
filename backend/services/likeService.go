package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"realTimeForum/backend/controllers"
	"realTimeForum/backend/models"
)

var like models.Like

func LikeWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade connection:", err)
		return
	}
	defer func() {
		conn.Close()
		removeConnection(conn)
	}()

	addConnection(conn)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		err = json.Unmarshal(message, &like)
		if err != nil {
			fmt.Println("Error unmarshaling data:", err)
			return
		}

		posts, Errs := controllers.HandleLikedata(w, like, r)

		ServicePosts.Posts = posts
		ServicePosts.Errors = Errs
	
		newPostResponse := map[string]interface{}{
			"type":   "new_post",
			"Posts":  posts,
			"Errors": Errs,
		}

		response, err := json.Marshal(newPostResponse)
		if err != nil {
			fmt.Println("Error marshaling response:", err)
			return
		}
	
		if Errs.Status {
			err = conn.WriteMessage(messageType, response)
			if err != nil {
				fmt.Println("Error writing error message:", err)
			}
		} else {
			broadcastMessage(messageType, response)
		}
	}
}
