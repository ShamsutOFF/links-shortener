package auth

import (
	"encoding/json"
	"fmt"
	"links-shortener/configs"
	"log"
	"net/http"
)

type AuthHandlerDeps struct {
	*configs.Config
}

type AuthHandler struct {
	*configs.Config
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println(handler.Config.Auth.Secret)
		fmt.Println("Login")
		res := LoginResponse{
			Token: "123",
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(200)
		err := json.NewEncoder(writer).Encode(res)
		if err != nil {
			log.Println(err, "error encode response")
		}
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("Register")
	}
}
