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
		http.SetCookie(w, &http.Cookie{
			Name:  "message",
			Value: "you are logged out",
		})
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	username, err := cfg.validateJwt(cookie.Value)
	if err != nil {
		http.SetCookie(w, &http.Cookie{
			Name:  "message",
			Value: "you are logged out",
		})
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	// at this stage user is logged in, check for any messages cookie, if not found then render username is loggin
	var d flashMsg

	msgCookie, err := r.Cookie("message")
	log.Println(msgCookie)
	if err != nil || msgCookie.Value == "" {
		d = flashMsg{
			Msg:   username + " is logged in",
			Color: "success",
		}
	} else {
		d = flashMsg{
			Msg:   msgCookie.Value,
			Color: "failure",
		}
	}

	cfg.tmpl.home.Execute(w, d)
}
func (cfg *Config) loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		cookie, err := r.Cookie("message")
		if err != nil {
			cfg.tmpl.login.Execute(w, nil)
		} else {
			e := flashMsg{
				Msg:   cookie.Value,
				Color: "failure",
			}
			cfg.tmpl.login.Execute(w, e)
		}
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
		removeCookie(w, "message")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (cfg *Config) joinRoomHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		roomCodeString := r.FormValue("roomCode")
		if roomCodeString == "" {
			log.Println("roomCode was empty")
			setCookie(w, "message", "roomCode field empty")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		roomCodeNumerical, _ := strconv.Atoi(roomCodeString)

		_, err := cfg.db.getRoomByCode(roomCodeNumerical)
		if err != nil {
			log.Println("roomCode is not correct")
			setCookie(w, "message", "incorrect roomCode")
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
	}
}

func (cfg *Config) createRoomHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// create the room and add the user into it
	}

}

// func broadcast(room Room, mt int, msg []byte) {
// 	// for _,con := range room.ClientList{
// 	// 	con.WriteMessage(mt, []byte(msg))
// 	// }
// }

func removeCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:   name,
		Value:  "",
		MaxAge: -1,
	})
}

func setCookie(w http.ResponseWriter, name string, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:   name,
		Value:  value,
		Domain: "/",
	})
}
