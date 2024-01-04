package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type flashMsg struct {
	Msg   string
	Color string
}

func (cfg *Config) wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("websocket connection established...")
	conn, err := cfg.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error in upgrading : ", err)
	}
	defer conn.Close()
	// room.ClientList = append(room.ClientList, conn)
	// infinite read loop
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read failed: ", err)
		}

		if string(message) == "the connection has opened" {

			log.Println("message recieved : ", string(message))
			conn.WriteMessage(mt, message)
		} else {
			log.Println("broadcasting message : ", string(message))
			// broadcast(room, mt, message)
		}
	}

}

func (cfg *Config) homeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		SetFlash(w, "errorMessage", "you need to login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username, err := cfg.validateJwt(cookie.Value)
	if err != nil {
		SetFlash(w, "errorMessage", "you need to login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	var d flashMsg
	msg := GetFlash(w, r, "errorMessage")
	if msg == "" {
		d = flashMsg{
			Msg:   username + " is logged in",
			Color: "success",
		}
	} else {
		d = flashMsg{
			Msg:   msg,
			Color: "failure",
		}
	}
	cfg.tmpl.home.Execute(w, d)
}

func (cfg *Config) signupHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		cfg.tmpl.signup.Execute(w, nil)
	case "POST":
		username := r.FormValue("username")
		displayname := r.FormValue("displayName")
		password := r.FormValue("password")
	}
}

func (cfg *Config) loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		msg := GetFlash(w, r, "errorMessage")
		if msg == "" {
			cfg.tmpl.login.Execute(w, nil)
			return
		}
		d := flashMsg{
			Msg:   msg,
			Color: "failure",
		}
		cfg.tmpl.login.Execute(w, d)
	case "POST":
		username := r.FormValue("username")
		password := r.FormValue("password")
		user, err := cfg.db.getUserByUserName(username)

		isCorrectPassword := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

		if err != nil || isCorrectPassword != nil {
			e := flashMsg{
				Msg:   "invalid credentials",
				Color: "failure",
			}
			cfg.tmpl.login.Execute(w, e)
			return
		}

		tokenString, err := cfg.generateJwt(username)
		if err != nil {
			e := flashMsg{
				Msg:   err.Error(),
				Color: "failure",
			}
			cfg.tmpl.login.Execute(w, e)
		}
		http.SetCookie(w, &http.Cookie{
			Name:    "jwt",
			Value:   tokenString,
			Expires: time.Now().Add(time.Minute * 2),
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (cfg *Config) joinRoomHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		roomCodeString := r.FormValue("roomCode")
		roomCodeNumerical, _ := strconv.Atoi(roomCodeString)

		_, err := cfg.db.getRoomByCode(roomCodeNumerical)
		if err != nil {
			SetFlash(w, "errorMessage", "room-code does not exist")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
}

func (cfg *Config) createRoomHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// create the room and add the user into it
		roomName := r.FormValue("roomName")
		cfg.db.createRoom(roomName)
	}

}

// func broadcast(room Room, mt int, msg []byte) {
// 	// for _,con := range room.ClientList{
// 	// 	con.WriteMessage(mt, []byte(msg))
// 	// }
// }
