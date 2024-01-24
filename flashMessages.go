package main

import (
	"net/http"
	"time"
)

// Sets flash
func SetFlash(w http.ResponseWriter, name string, value string) {
	c := http.Cookie{
		Name:  name,
		Value: value,
		Path:  "/",
	}
	http.SetCookie(w, &c)
}

// Gets flash
func GetFlash(w http.ResponseWriter, r *http.Request, name string) string {
	c, err := r.Cookie(name)
	if err != nil {
		return ""
	}
	valueToReturn := c.Value

	removeCookie := http.Cookie{
		Name:    name,
		MaxAge:  -1,
		Expires: time.Unix(1, 0),
	}
	http.SetCookie(w, &removeCookie)

	return valueToReturn
}
