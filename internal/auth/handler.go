package auth

import (
	"fmt"
	"links-shortener/pkg/req"
	"net/http"

	"links-shortener/configs"
	"links-shortener/pkg/res"
)

type AuthHandlerDeps struct {
	*configs.Config
	*AuthService
}

type AuthHandler struct {
	*configs.Config
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDeps) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("Login()")
		body, err := req.HandleBody[LoginRequest](&writer, request)
		if err != nil {
			return
		}
		fmt.Println(body)
		data := LoginResponse{
			Token: "123",
		}
		res.JsonResp(writer, data, http.StatusOK)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&writer, request)
		if err != nil {
			return
		}
		handler.AuthService.Register(body.Email, body.Password, body.Name)
	}
}
