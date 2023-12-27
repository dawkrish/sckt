package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct{
	db *databaseConfig
	tmpl *templateConfig
	upgrader websocket.Upgrader
	JWT_SECRET []byte
}

type templateConfig  struct {
	home *template.Template
	login *template.Template
}
type databaseConfig struct {
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

var cfg Config

func main() {
	cfg , err := initalizeCfg()
	if err != nil {
		log.Fatal("error initalizing cfg : " + err.Error())
	} 
	defer cfg.db.mongoClient.Disconnect(context.TODO())

	r := http.NewServeMux()
	fs := http.FileServer(http.Dir("./static"))

	r.Handle("/static/", http.StripPrefix("/static/", fs))
	r.HandleFunc("/login", cfg.loginHandler)
	r.HandleFunc("/",cfg.homeHandler)

	r.HandleFunc("/ws", cfg.wsHandler)

	log.Println("listening on http://localhost:8080")
	http.ListenAndServe("localhost:8080", r)
}


func initalizeCfg() (Config, error){
	homeTemplate, err := template.ParseFiles("./templates/home.html")
	if err != nil {
		return Config{},errors.New("error parsing home.html : " + err.Error())
	}
	loginTemplate, err := template.ParseFiles("./templates/login.html")
	if err != nil {
		return Config{},errors.New("error parsing login.html : " + err.Error())
	}
	tmpl := templateConfig{
		home : homeTemplate,
		login: loginTemplate,
	}	
	upgrader := websocket.Upgrader{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}
	if err = godotenv.Load(); err != nil{
		return Config{},errors.New("error .env file not found : "+ err.Error())
	}
	JWT_SECRET := os.Getenv("JWT_SECRET")
	if JWT_SECRET == ""{
		return Config{}, errors.New("error JWT_SECRET not found in .env : " + err.Error())
	}
	MONGO_URI := os.Getenv("MONGO_URI")
	if MONGO_URI == "" {
		return Config{}, errors.New("error MONGO_URI not found in .env : " + err.Error())
	}

	mongoClient , err := mongo.Connect(context.TODO(),options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		return Config{}, errors.New("error could not connect to mongoDB : " + err.Error())
	}

	dbCfg := databaseConfig{
		mongoClient: mongoClient,
		userColl: mongoClient.Database("sckt").Collection("user"),
		roomColl: mongoClient.Database("sckt").Collection("room"),
		messageColl: mongoClient.Database("sckt").Collection("message"),
	}
	
	cfg := Config{
		db : &dbCfg,
		tmpl: &tmpl,
		upgrader: upgrader,
		JWT_SECRET: []byte(JWT_SECRET),
	}
	return cfg,nil
}
