package main

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	// "github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

type flashMsg struct {
	Msg   string
	Color string
}

func (cfg *Config) isClientExists(newClient Client) (Client, int) {
	for _, v := range cfg.Clients {
		if v.RoomCode == newClient.RoomCode && v.Username == newClient.Username {
			return v, 1
		}
	}
	return Client{}, -1
}

func (cfg *Config) wsHandler(w http.ResponseWriter, r *http.Request) {
	username, err := cfg.middlewareJwt(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	code, _ := strconv.Atoi(chi.URLParam(r, "code"))
	client := Client{
		Username: username,
		RoomCode: code,
	}
	pseudoClient, existCode := cfg.isClientExists(client)
	if existCode == 1 {
		log.Println("this client is already in the room")
		pseudoClient.Conn.Close()
		w.Write([]byte("one window is already open for this chat..."))
		return
	}
	conn, err := cfg.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("error in upgrading : ", err)
	}
	client.Conn = conn
	cfg.Clients = append(cfg.Clients, client)
	log.Println(cfg.Clients)
	log.Printf("websocket-connection-%v-established-by-%v\n", code, username)
	for {
		mt, message, err := conn.ReadMessage()
		log.Println("message recieved : ", string(message))
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
			customClients := []Client{}
			for _, v := range cfg.Clients {
				if v.RoomCode == code {
					customClients = append(customClients, v)
				}
			}
			broadcast(mt, message, customClients)
		}
	}

}

func broadcast(mt int, msg []byte, customClients []Client) {
	log.Println("the clients i will broadcast to : ", customClients)
	for _, v := range customClients {
		v.Conn.WriteMessage(mt, msg)
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
	username, err := cfg.middlewareJwt(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
	user, _ := cfg.db.getUserByUserName(username)
	log.Println("the user that will join room... ", user)
	roomCodeString := r.FormValue("roomCode")
	roomCodeNumerical, _ := strconv.Atoi(roomCodeString)

	room, err := cfg.db.getRoomByCode(roomCodeNumerical)
	if err != nil {
		SetFlash(w, "errorMessage", "room-code does not exist")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	filter := bson.D{{Key: "code", Value: room.Code}}
	update := bson.D{{Key: "$push", Value: bson.D{{Key: "users", Value: user}}}}

	result, err := cfg.db.roomColl.UpdateOne(context.TODO(), filter, update)
	log.Println("error in updating room's user's ", err)
	log.Printf("Matched Count : %v \n Updated Count : %v\n", result.MatchedCount, result.ModifiedCount)

}

func (cfg *Config) createRoomHandler(w http.ResponseWriter, r *http.Request) {
	username, err := cfg.middlewareJwt(w, r)
	if err != nil {
		SetFlash(w, "errorMessage", "you need to login")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	roomName := r.FormValue("roomName")
	user, _ := cfg.db.getUserByUserName(username)
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
