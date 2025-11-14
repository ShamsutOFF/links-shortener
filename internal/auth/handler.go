package auth

import (
	"links-shortener/pkg/jwt"
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
			return
		}
		token, err := jwt.NewJWT(handler.Config.Auth.Secret).GenerateToken(login)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		data := LoginResponse{
			Token: token,
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
		login, err := handler.AuthService.Register(body.Email, body.Password, body.Name)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		token, err := jwt.NewJWT(handler.Config.Auth.Secret).GenerateToken(login)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}
		data := RegisterResponse{
			Token: token,
		}
		res.JsonResp(writer, data, http.StatusOK)
	}
}
