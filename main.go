package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
)

type DatabaseConfig struct{
	mongoClient *mongo.Client
	userColl *mongo.Collection

}

type Room struct{
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
	getAllUsers()
	r := mux.NewRouter()

	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/ws", wsHandler)

	log.Println("listening on http://localhost:8080")
	http.ListenAndServe("localhost:8080", r)
}
