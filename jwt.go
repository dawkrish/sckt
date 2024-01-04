package main

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
)

func (cfg *Config) generateJwt(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(2 * time.Minute).Unix(),
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
