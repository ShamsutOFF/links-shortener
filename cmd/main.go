package main

import (
	"fmt"
	"links-shortener/internal/auth"
	"net/http"
)

func main() {
	//_ := configs.LoadConfig()
	router := http.NewServeMux()
	auth.NewAuthHandler(router)

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
