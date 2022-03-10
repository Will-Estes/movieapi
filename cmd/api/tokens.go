package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"movieapi/models"
	"net/http"
	"time"
)

var validUser = models.User{
	ID:       1,
	Email:    "test@test.com",
	Password: "$2a$12$YGrq1IXO76oAuSybFAp4JOi4kSAAR3lC33YArhHlkVPBD73z4jEVi",
}

type Credentials struct {
	Username string `json:"email"`
	Password string `json:"password"`
}

func (app *application) SignIn(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		app.errorJSON(w, errors.New("decode error"))
		return
	}

	// This would be a DB call in a real app
	hashedPassword := validUser.Password

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password))
	if err != nil {
		app.errorJSON(w, errors.New("unauthorized"))
		return
	}

	var claims jwt.StandardClaims
	claims.Subject = fmt.Sprint(validUser.ID)
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(24 * time.Hour).Unix()
	claims.Issuer = "mydomain.com"
	claims.Audience = "mydomain.com"

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(app.config.jwt.secret))
	if err != nil {
		app.errorJSON(w, errors.New("error signing"))
		return
	}
	err = app.writeJSON(w, http.StatusOK, tokenString, "response")
	if err != nil {
		app.errorJSON(w, errors.New("error writing response"))
		return
	}
}
