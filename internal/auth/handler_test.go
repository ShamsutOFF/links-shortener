package auth_test

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"links-shortener/configs"
	"links-shortener/internal/auth"
	"links-shortener/internal/user"
	"links-shortener/pkg/db"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	TestEmail = "aaa23@ya.ru"
)

func bootstrap() (*auth.AuthHandler, sqlmock.Sqlmock, error) {
	database, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn: database,
	}))
	if err != nil {
		return nil, nil, err
	}
	userRepo := user.NewUserRepository(&db.Db{
		DB: gormDb,
	})
	handler := auth.AuthHandler{
		Config: &configs.Config{
			Auth: configs.AuthConfig{
				Secret: "secret",
			},
		},
		AuthService: auth.NewAuthService(userRepo),
	}
	return &handler, mock, nil
}

func TestAuthHandler_LoginSuccess(t *testing.T) {
	handler, mock, err := bootstrap()
	rows := sqlmock.NewRows([]string{"email", "password"}).
		AddRow(TestEmail, "$2a$10$/bGD7PqJXPKUnFa9oSz1Yuh8ksf.1G8zJ9nSxFsFyna631KEqYX16")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	if err != nil {
		t.Fatal(err)
		return
	}
	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    TestEmail,
		Password: "passs",
	})
	reader := bytes.NewReader(data)

	wr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/auth/login", reader)
	handler.Login()(wr, req)
	if wr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, wr.Code)
	}
}
