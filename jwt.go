package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

func (cfg *Config) generateJwt(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(3 * time.Minute).Unix(),
		Subject:   username,
	})
	tokenString, err := token.SignedString(cfg.JWT_SECRET)
	return tokenString, err
}

func (cfg *Config) validateJwt(tokenString string) (string, error) {
	customClaims := jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &customClaims, func(t *jwt.Token) (interface{}, error) {
		return cfg.JWT_SECRET, nil
	})
	if err != nil {
		return "", errors.New("error parsing token : " + err.Error())
	}
	if !token.Valid {
		return "", errors.New("error token not valid : " + err.Error())
	}
	username := customClaims.Subject
	return username, nil
}

func (cfg *Config) middlewareJwt(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		SetFlash(w, "errorMessage", "you need to login")
		return "", errors.New("you need to login")
	}
	username, err := cfg.validateJwt(cookie.Value)
	if err != nil {
		SetFlash(w, "errorMessage", "you need to login")
		return "", errors.New("you need to login")
	}
	return username, nil
}
