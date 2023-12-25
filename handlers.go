package main

import (
	"html/template"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	// "github.com/gorilla/websocket"
	// "encoding/json"
	// "github.com/golang-jwt/jwt"
)

type errToSend struct {
	Emsg string
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("websocket connection established...")
	conn, err := upgrader.Upgrade(w, r, nil)
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
			broadcast(room, mt, message)
		}
	}

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	// first check if there is a user logged in or not, if not then then send to login page or else display home page
	// isUserLoggedIn := 
	t, _ := template.ParseFiles("templates/home.html")
	type dataToSend struct {
		Title string
	}
	d := dataToSend{
		Title: "Home Page",
	}
	t.Execute(w, d)
}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		t, _ := template.ParseFiles("templates/login.html")
		t.Execute(w, nil)
	case "POST":
		username := r.FormValue("username")
		password := r.FormValue("password")
		user, err := getUserByUserName(username)
		
		isCorrectPassword := bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(password))

		if err != nil || isCorrectPassword != nil {
			e := errToSend{
				Emsg: "invalid credentials",
			}
			tmplCfg.login.Execute(w, e)
			return
		}
		// create a jwt and assign it to a cookie....

		tokenString,err :=generateJwt(username)
		if err != nil {
			e := errToSend{
				Emsg: err.Error(),
			}
			tmplCfg.login.Execute(w,e)
		}
		http.SetCookie(w,&http.Cookie{
			Name: "jwt",
			Value: tokenString,
		})
		http.Redirect(w,r,"/",http.StatusSeeOther)
	}
}
func broadcast(room Room, mt int, msg []byte) {
	// for _,con := range room.ClientList{
	// 	con.WriteMessage(mt, []byte(msg))
	// }
}
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDM1MzAxNDcsInN1YiI6InZhbnNoIn0.JtqUL71NhyGuDKAtKK_I9zFYkcb3_-sCZmdSYYtFy2k