package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
)

type DatabaseConfig struct{
	mongoClient *mongo.Client
	userColl *mongo.Collection
	roomColl *mongo.Collection
	messageColl *mongo.Collection

}

type ClientRoom struct{
	Code int
	Name string
	ClientList []*websocket.Conn
}

var (
	room Room
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	dbCfg *DatabaseConfig = new(DatabaseConfig)
)

func main() {
	conncectToDB()
	defer dbCfg.mongoClient.Disconnect(context.TODO())
	// createUser("vansh","c@c.com", "root")
	// createUser("anant","d@d.com", "root")
	// getAllUsers()
	// deleteAllUsers()
	// createRoom(3478)
	r := http.NewServeMux()	
	fs := http.FileServer(http.Dir("./static"))

	r.Handle("/static/",http.StripPrefix("/static/",fs))
	r.HandleFunc("/login",loginHandler)

	r.HandleFunc("/ws", wsHandler)

	log.Println("listening on http://localhost:8080")
	http.ListenAndServe("localhost:8080", r)
}
