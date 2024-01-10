package main

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	// "github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

type flashMsg struct {
	Msg   string
	Color string
}

func ifClientExists(clientList []Client, newClient Client) int {
	for in, val := range clientList {
		if val.Username == newClient.Username {
			return in
		}
	}
	return -1
}

func (cfg *Config) isClientRoomExists(newCode int) (ClientRoom, error) {
	for _, val := range cfg.ClientRooms {
		if val.Code == newCode {
			return val, nil
		}
	}
	return ClientRoom{}, errors.New("no clientRoom exists with this code")
}

func (cfg *Config) wsHandler(w http.ResponseWriter, r *http.Request) {
	username, err := cfg.middlewareJwt(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	code, _ := strconv.Atoi(chi.URLParam(r, "code"))
	conn, err := cfg.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error in upgrading : ", err)
	}
	client := Client{
		Username: username,
		Conn:     conn,
	}
	clientRoom := ClientRoom{
		Code: code,
	}
	// we need to check that this client room code exists already or not
	// if cfg.isClientRoomExists(code) == -1 {
	// 	clientRoom := ClientRoom{
	// 		Code: code,
	// 	}
	// }

	// if ifClientExists(clientRoom.ClientList)
	// we need to check that the client room we are adding, should not have
	clientRoom.ClientList = append(clientRoom.ClientList, client)
	cfg.ClientRooms = append(cfg.ClientRooms, clientRoom)
	log.Printf("websocket-connection-%v-established-by-%v\n", code, username)
	// defer conn.Close()
	// infinite read loop
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read failed: ", err)
			conn.Close()
			break
		}

		if string(message) == "the connection has opened" {

			log.Println("message recieved : ", string(message))
			conn.WriteMessage(mt, message)
		} else if string(message) == "the connection has closed" {
			log.Println("message recieved : ", string(message))
			conn.Close()
		} else {
			log.Println("broadcasting message : ", string(message))
			// broadcast(room, mt, message)
		}
	}

}

func (cfg *Config) homeHandler(w http.ResponseWriter, r *http.Request) {
	username, err := cfg.middlewareJwt(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	var data struct {
		Msg   string
		Color string
		Rooms []Room
	}
	msg := GetFlash(w, r, "errorMessage")
	if msg == "" {
		data.Msg = username + " is logged in"
		data.Color = "success"
	} else {
		data.Msg = msg
		data.Color = "failure"
	}
	rooms, _ := cfg.db.getAllRooms(username)
	data.Rooms = rooms

	cfg.tmpl.home.Execute(w, data)

}

func (cfg *Config) GetSignupHandler(w http.ResponseWriter, r *http.Request) {
	cfg.tmpl.signup.Execute(w, nil)
}

func (cfg *Config) PostSignupHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	displayname := r.FormValue("displayName")
	email := r.FormValue("email")
	password := r.FormValue("password")

	_, ifUsernameExists := cfg.db.getUserByUserName(username)
	if ifUsernameExists == nil {
		d := flashMsg{
			Msg:   "username already exists",
			Color: "failure",
		}
		cfg.tmpl.signup.Execute(w, d)
		return
	}
	_, ifEmailExists := cfg.db.getUserByEmail(email)
	if ifEmailExists == nil {
		d := flashMsg{
			Msg:   "email already exists",
			Color: "failure",
		}
		cfg.tmpl.signup.Execute(w, d)
		return
	}
	newUser, err := cfg.db.createUser(username, displayname, email, password)
	if err != nil {
		log.Println("error creating new user : " + err.Error())
	}
	log.Println(newUser)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (cfg *Config) GetLoginHandler(w http.ResponseWriter, r *http.Request) {
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

}
func (cfg *Config) PostLoginHandler(w http.ResponseWriter, r *http.Request) {
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

func (cfg *Config) joinRoomHandler(w http.ResponseWriter, r *http.Request) {
	roomCodeString := r.FormValue("roomCode")
	roomCodeNumerical, _ := strconv.Atoi(roomCodeString)

	_, err := cfg.db.getRoomByCode(roomCodeNumerical)
	if err != nil {
		SetFlash(w, "errorMessage", "room-code does not exist")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}

func (cfg *Config) createRoomHandler(w http.ResponseWriter, r *http.Request) {
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
	roomName := r.FormValue("roomName")
	user, err := cfg.db.getUserByUserName(username)
	room, err := cfg.db.createRoom(user, roomName)
	if err != nil {
		log.Println("error creating room : " + err.Error())
	}
	log.Println("room created : ", room)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (cfg *Config) chatHandler(w http.ResponseWriter, r *http.Request) {
	username, err := cfg.middlewareJwt(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	code, _ := strconv.Atoi(chi.URLParam(r, "code"))

	var data struct {
		Code     int
		Username string
	}
	data.Code = code
	data.Username = username
	cfg.tmpl.chat.Execute(w, data)
}

// func broadcast(room Room, mt int, msg []byte) {
// 	// for _,con := range room.ClientList{
// 	// 	con.WriteMessage(mt, []byte(msg))
// 	// }
// }
