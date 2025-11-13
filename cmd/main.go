package main

import (
	"fmt"
	"links-shortener/configs"
	"links-shortener/internal/auth"
	"links-shortener/internal/link"
	"links-shortener/pkg/db"
	"links-shortener/pkg/middleware"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	newDb := db.NewDb(conf)
	router := http.NewServeMux()

	// Repositories
	linkRepository := link.NewLinkRepository(newDb)

	// Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{Config: conf})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
	})

	// Middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)
	server := http.Server{
		Addr:    ":7777",
		Handler: stack(router),
	}

	fmt.Println("Server listening on port 7777")
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Server listening error")
		return
	}
}
