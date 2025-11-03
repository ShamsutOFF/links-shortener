package main

import (
	"fmt"
	"links-shortener/configs"
	"links-shortener/internal/auth"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	router := http.NewServeMux()
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{Config: conf})

	server := http.Server{
		Addr:    ":7777",
		Handler: router,
	}

	fmt.Println("Server listening on port 7777")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Server listening error")
		return
	}
}
