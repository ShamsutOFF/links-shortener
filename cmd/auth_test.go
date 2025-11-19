package main

import (
	"bytes"
	"encoding/json"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io"
	"links-shortener/internal/auth"
	"links-shortener/internal/user"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

const (
	TestEmail = "aaa23@ya.ru"
)

func initDb() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(os.Getenv("DSN")), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

func initData(db *gorm.DB) {
	db.Create(&user.User{
		Email:    TestEmail,
		Name:     "Вася Тестов",
		Password: "$2a$10$/bGD7PqJXPKUnFa9oSz1Yuh8ksf.1G8zJ9nSxFsFyna631KEqYX16",
	})
}

func removeData(db *gorm.DB) {
	db.Unscoped().
		Where("email = ?", TestEmail).
		Delete(&user.User{})
}

func TestLoginSuccess(t *testing.T) {
	// Prepare
	db := initDb()
	initData(db)

	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    TestEmail,
		Password: "passs",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	var resData auth.LoginResponse
	err = json.Unmarshal(body, &resData)
	if err != nil {
		t.Fatal(err)
	}
	if resData.Token == "" {
		t.Fatal("expected token, got empty")
	}
	removeData(db)
}

func TestLoginFail(t *testing.T) {
	// Prepare
	db := initDb()
	initData(db)

	ts := httptest.NewServer(App())
	defer ts.Close()

	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "notexist@ya.ru",
		Password: "notexist",
	})

	res, err := http.Post(ts.URL+"/auth/login", "application/json", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected status code %d, got %d", http.StatusUnauthorized, res.StatusCode)
	}
	removeData(db)
}
