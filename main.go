package main

import (
	"context"
	"log"
	"net/http"
	"text/template"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/mongo"
)

type DatabaseConfig struct {
	mongoClient *mongo.Client
	userColl    *mongo.Collection
	roomColl    *mongo.Collection
	messageColl *mongo.Collection
}

type ClientRoom struct {
	Code       int
	Name       string
	ClientList []*websocket.Conn
}

var (
	room     Room
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	dbCfg   *DatabaseConfig = new(DatabaseConfig)
	tmplCfg struct {
		login *template.Template
	}
)

func main() {
	createAllTemplates()
	conncectToDB()
	defer dbCfg.mongoClient.Disconnect(context.TODO())
	// createUser("vansh","Vansh Aggarwal","c@c.com", "root")
	// createUser("anant","Anant Gupta","d@d.com", "root")
	// getAllUsers()
	// deleteAllUsers()
	// createRoom(3478)
	r := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))

	r.Handle("/static/", http.StripPrefix("/static/", fs))
	r.HandleFunc("/login", loginHandler)
	r.HandleFunc("/",homeHandler)

	r.HandleFunc("/ws", wsHandler)

	log.Println("listening on http://localhost:8080")
	http.ListenAndServe("localhost:8080", r)
}

func createAllTemplates() {
	tmplCfg.login, _ = template.ParseFiles("./templates/login.html")
}
