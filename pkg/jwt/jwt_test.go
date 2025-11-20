package jwt_test

import (
	"links-shortener/pkg/jwt"
	"testing"
)

func TestJWT_GenerateToken(t *testing.T) {
	const email = "test@test.com"
	jwtService := jwt.NewJWT("–,3;8ZyVY09D5A,GOw?%x79iKn/O^_ku")
	token, err := jwtService.GenerateToken(jwt.JWTData{
		Email: email,
	})
	if err != nil {
		t.Fatal(err)
	}
	isValid, data := jwtService.Parse(token)
	if !isValid {
		t.Fatal("Token is not valid")
	}
	if data.Email != email {
		t.Fatalf("Email %s is not equal to %s", data.Email, email)
	}
}
