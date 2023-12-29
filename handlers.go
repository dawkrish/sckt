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

type flashMsg struct {
	Msg string
	Color string
}


func (cfg *Config)wsHandler(w http.ResponseWriter, r *http.Request) {
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

func (cfg *Config)homeHandler(w http.ResponseWriter, r *http.Request) {
	// first check if there is a user logged in or not, if not then then send to login page or else display home page
	cookie, err := r.Cookie("jwt")
	if err != nil {
		log.Println("error cookie not found : " + err.Error() )
	}
	username,err := cfg.validateJwt(cookie.Value) 	
	if err != nil {
		log.Println(err.Error())
		removeCookie(w,"jwt")
		http.SetCookie(w,&http.Cookie{
			Name: "message",
			Value: "you are logged out",
		})
		http.Redirect(w,r,"/login",http.StatusSeeOther)
	}
	removeCookie(w,"message")
	
	d := flashMsg{
		Msg: username + " is logged in",
		Color: "success",
	}	
	cfg.tmpl.home.Execute(w,d)
}
func (cfg *Config)loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		cookie, err := r.Cookie("message")
		if err != nil {
			log.Println("error message cookie not found : " + err.Error())
		}
		log.Println("cookie value : ",cookie.Value)
		e := flashMsg{
			Msg : cookie.Value,
			Color: "failure",
		}
		t, _ := template.ParseFiles("templates/login.html")
		t.Execute(w, e)
	case "POST":
		username := r.FormValue("username")
		password := r.FormValue("password")
		user, err := cfg.db.getUserByUserName(username)
		
		isCorrectPassword := bcrypt.CompareHashAndPassword([]byte(user.Password),[]byte(password))

		if err != nil || isCorrectPassword != nil {
			e := flashMsg{
				Msg: "invalid credentials",
				Color : "failure",
			}
			cfg.tmpl.login.Execute(w, e)
			return
		}

		tokenString,err := cfg.generateJwt(username)
		if err != nil {
			e := flashMsg{
				Msg: err.Error(),
				Color : "failure",
			}
			cfg.tmpl.login.Execute(w,e)
		}
		http.SetCookie(w,&http.Cookie{
			Name: "jwt",
			Value: tokenString,
		})
		http.Redirect(w,r,"/",http.StatusSeeOther)
	}
}
// func broadcast(room Room, mt int, msg []byte) {
// 	// for _,con := range room.ClientList{
// 	// 	con.WriteMessage(mt, []byte(msg))
// 	// }
// }

func removeCookie(w http.ResponseWriter, name string){
	http.SetCookie(w,&http.Cookie{
		Name: name,
		Value: "",
		MaxAge: -1,
	})
}