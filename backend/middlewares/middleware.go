package middlewares

import (
	"fmt"
	"net/http"
	"realTimeForum/backend/routes"
	// "realTimeForum/backend"
)

func ErrorMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		foundpath, foundmethod := false, false
		for _, route := range routes.Routes {
			if route.Path == r.URL.Path {
				foundpath = true
				for _, method := range route.Method {
					if r.Method == method {
						foundmethod = true
					}
				}
			}
		}
		if !foundpath {
			w.WriteHeader(404)
			// NotFound
			fmt.Println(r.URL.Path)
			fmt.Println("Not Found")
			return
		}
		if !foundmethod {
			w.WriteHeader(405)
			//Method Not Allowed
			fmt.Println("Method Not Allowed")
			return
		}
		next.ServeHTTP(w, r)
	})
}
