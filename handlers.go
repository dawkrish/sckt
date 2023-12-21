package main

import (
	"html/template"
	"log"
	"net/http"
	// "github.com/gorilla/websocket"
	// "encoding/json"
)

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
	switch r.Method{
	case "GET":
		t, _ := template.ParseFiles("templates/login.html")
		t.Execute(w, nil)
	case "POST":
		username := r.FormValue("username")
		password := r.FormValue("password")
		log.Println(username,password)
		w.Write([]byte("/login POST"))
}
	}
func POSTloginHandler(w http.ResponseWriter, r *http.Request) {

}
func broadcast(room Room, mt int, msg []byte) {
	// for _,con := range room.ClientList{
	// 	con.WriteMessage(mt, []byte(msg))
	// }
}
