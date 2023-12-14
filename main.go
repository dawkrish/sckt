package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var ClientList []*websocket.Conn

// type Client struct{
// 	conn *websocket.Conn
// }

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

func main(){
	r := mux.NewRouter()
	
	r.HandleFunc("/",homeHandler)
	r.HandleFunc("/ws",wsHandler)	
	

	log.Println("listening on http://localhost:8080")
	http.ListenAndServe("localhost:8080",r)
}
