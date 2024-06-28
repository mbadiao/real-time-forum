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

var post models.Post
var ServicePosts models.ServicePost
var activeConnectionsMutex sync.Mutex
var activeConnections = make(map[*websocket.Conn]bool)
var ErrorState = &struct {
	Status  bool
	Code    int
	Message string
}{}

func PostWebsocket(w http.ResponseWriter, r *http.Request) {
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

		var req map[string]interface{}
		err = json.Unmarshal(message, &req)
		if err != nil {
			fmt.Println("Error unmarshaling data:", err)
			continue
		}

		if req["type"] == "new_post" {
			InsertPostSocket(w, r, messageType, message, conn)
		} else if req["type"] == "get_all_posts" {
			GetPostSocket(messageType, conn)
		}
	}
}

func InsertPostSocket(w http.ResponseWriter, r *http.Request, messageType int, message []byte, conn *websocket.Conn) {
	err := json.Unmarshal(message, &post)
	if err != nil {
		fmt.Println("Error unmarshaling data:", err)
		return
	}

	posts, Errs := controllers.HandlePostData(w, post, r)
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

func GetPostSocket(messageType int, conn *websocket.Conn) {
	posts, err := controllers.GetAllPosts()
	if err != nil {
		fmt.Println("Error getting all posts:", err)
		GetErrorS("InternalServerError", 500)
	}
	response := map[string]interface{}{
		"type":   "all_posts",
		"Posts":  posts,
		"Errors": ErrorState,
	}
	responseData, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshaling response:", err)
		GetErrorS("InternalServerError", 500)
		return
	}
	if ErrorState.Status {
		err = conn.WriteMessage(messageType, responseData)
		if err != nil {
			fmt.Println("Error writing error message:", err)
		}
	} else {
		broadcastMessage(messageType, responseData)
	}

}

func addConnection(conn *websocket.Conn) {
	activeConnectionsMutex.Lock()
	defer activeConnectionsMutex.Unlock()
	activeConnections[conn] = true
}

func removeConnection(conn *websocket.Conn) {
	activeConnectionsMutex.Lock()
	defer activeConnectionsMutex.Unlock()
	delete(activeConnections, conn)
}

func broadcastMessage(messageType int, data []byte) {
	fmt.Println("broadcastMessage")
	activeConnectionsMutex.Lock()
	defer activeConnectionsMutex.Unlock()
	fmt.Println("ACtive", activeConnections)
	for conn := range activeConnections {
		err := conn.WriteMessage(messageType, data)
		if err != nil {
			fmt.Println("Error writing message to connection:", err)
			removeConnection(conn)
		}
	}
}

func GetErrorS(msg string, code int) {
	ErrorState.Status = true
	ErrorState.Code = code
	ErrorState.Message = msg
}
