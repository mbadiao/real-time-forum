package routes

import (
	"net/http"
	"realTimeForum/backend/controllers"
	"realTimeForum/backend/services"
)

var Port = ":8080"

type Error struct {
	Code    int
	Message string
}

var Err = map[int]Error{
	404: {
		http.StatusNotFound,
		http.StatusText(404),
	},
	500: {
		http.StatusInternalServerError,
		http.StatusText(500),
	},
	400: {
		http.StatusBadRequest,
		http.StatusText(400),
	},
	405: {
		http.StatusMethodNotAllowed,
		http.StatusText(http.StatusMethodNotAllowed),
	},
}

type Route struct {
	Path    string
	Handler http.HandlerFunc
	Method  []string
}

var Routes = []Route{
	{
		Path:    "/",
		Handler: controllers.HomeHandler,
		Method:  []string{"GET", "POST"},
	},
	{
		Path:    "/register",
		Handler: controllers.RegisterHandler,
		Method:  []string{"GET", "POST"},
	},
	{
		Path:    "/login",
		Handler: controllers.LoginHandler,
		Method:  []string{"GET", "POST"},
	},
	{
		Path:    "/messages",
		Handler: services.WsMessages,
		Method:  []string{"GET", "POST"},
	},
	{
		Path:    "/post",
		Handler: services.PostWebsocket,
		Method:  []string{"GET", "POST"},
	},
	{
		Path:    "/like",
		Handler: services.LikeWebsocket,
		Method:  []string{"GET", "POST"},
	},
	{
		Path:    "/comment",
		Handler: services.CommentWebsocket,
		Method:  []string{"GET", "POST"},
	},
	{
		Path:    "/userstatus",
		Handler: controllers.HandleOnlineUser,
		Method:  []string{"GET", "POST"},
	},
	{
		Path:    "/checksession",
		Handler: controllers.IsRightSession,
		Method:  []string{"GET", "POST"},
	},
	{
		Path:    "/notifications",
		Handler: services.Notifshandler,
		Method:  []string{"GET", "POST"},
	},
	{
		Path:    "/typing-progress",
		Handler: services.TypingHandler,
		Method:  []string{"GET", "POST"},
	},
}
