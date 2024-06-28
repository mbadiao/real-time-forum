package main

import (
	"fmt"
	"net/http"
	"os"
	"realTimeForum/backend/controllers"
	"realTimeForum/backend/middlewares"
	"realTimeForum/backend/routes"
	"realTimeForum/data"
	// "realTimeForum/backend/controllers"
)

var serves = []string{"static", "js"}

var db = data.CreateTable()

func main() {
	if len(os.Args) == 1 {
		for _, serve := range serves {
			fs := http.FileServer(http.Dir("./frontend/" + serve))
			http.Handle("/"+serve+"/", http.StripPrefix("/"+serve+"/", fs))
		}

		for _, route := range routes.Routes {
			http.Handle(route.Path, middlewares.ErrorMiddleware(route.Handler))
		}
		go controllers.HandleMessages()
		fmt.Println("Server running at:")
		fmt.Println("> Localhost:    \033[34mhttp://localhost" + routes.Port + "\033[0m")
		fmt.Println("> disconnect:   \033[31mpress Ctrl+C\033[0m")
		http.ListenAndServe(routes.Port, nil)
	}
}
