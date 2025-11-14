package auth

import (
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
		body, err := req.HandleBody[LoginRequest](&writer, request)
		if err != nil {
			return
		}
		login, err := handler.AuthService.Login(body.Email, body.Password)
		if err != nil {
			http.Error(writer, ErrWrongCredentials, http.StatusUnauthorized)
		}
		res.JsonResp(writer, login, http.StatusOK)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&writer, request)
		if err != nil {
			return
		}
		register, err := handler.AuthService.Register(body.Email, body.Password, body.Name)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		res.JsonResp(writer, register, http.StatusOK)
	}
}
