package main

import (
	"errors"
	"github.com/golang-jwt/jwt"
	"net/http"
	"strconv"
	"strings"
)

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		next.ServeHTTP(w, r)
	})
}

func (app *application) checkToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")
		//if authHeader === '' {
		//	// set an anonymous user?
		//}
		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 {
			app.errorJSON(w, errors.New("invalid auth header"))
			return
		}
		if headerParts[0] != "Bearer" {
			app.errorJSON(w, errors.New("unauthorized, no bearer"))
			return
		}
		token := headerParts[1]
		hmacSecret := []byte(app.config.jwt.secret)
		decodedToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
			return hmacSecret, nil
		})
		if err != nil {
			app.errorJSON(w, errors.New("error parsing token"), http.StatusForbidden)
			return
		}
		err = decodedToken.Claims.Valid()
		if err != nil {
			app.errorJSON(w, errors.New("unauthorized - token expired"), http.StatusForbidden)
			return
		}

		claims := decodedToken.Claims.(*jwt.StandardClaims)
		if !claims.VerifyAudience("mydomain.com", true) {
			app.errorJSON(w, errors.New("unauthorized - failed hmac check"), http.StatusForbidden)
			return
		}

		if !claims.VerifyIssuer("mydomain.com", true) {
			app.errorJSON(w, errors.New("unauthorized - invalid issuer"), http.StatusForbidden)
			return
		}

		userID, err := strconv.ParseInt(claims.Subject, 10, 64)
		if err != nil {
			app.errorJSON(w, errors.New("unauthorized"), http.StatusForbidden)
			return
		}

		app.logger.Println(userID)
		next.ServeHTTP(w, r)
	})
}
