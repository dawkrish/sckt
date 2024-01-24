package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	db         *databaseConfig
	tmpl       *templateConfig
	upgrader   websocket.Upgrader
	JWT_SECRET []byte
	Clients    []Client
}

type templateConfig struct {
	home   *template.Template
	login  *template.Template
	signup *template.Template
	chat   *template.Template
}

type databaseConfig struct {
	mongoClient *mongo.Client
	userColl    *mongo.Collection
	roomColl    *mongo.Collection
	messageColl *mongo.Collection
}

type Client struct {
	Conn     *websocket.Conn
	Username string
	RoomCode int
}

func main() {
	cfg, err := initalizeCfg()
	if err != nil {
		log.Fatal("error initalizing cfg : " + err.Error())
	}
	defer cfg.db.mongoClient.Disconnect(context.TODO())

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	fs := http.FileServer(http.Dir("./static"))
	r.Handle("/static/*", http.StripPrefix("/static/", fs))

	r.Get("/", cfg.homeHandler)

	r.Get("/login", cfg.GetLoginHandler)
	r.Post("/login", cfg.PostLoginHandler)
	r.Get("/signup", cfg.GetSignupHandler)
	r.Post("/signup", cfg.PostSignupHandler)
	r.Post("/room/create", cfg.createRoomHandler)
	r.Post("/room/join", cfg.joinRoomHandler)
	r.Get("/chat/{code:[0-9]+}", cfg.chatHandler)

	r.Get("/ws/{code:[0-9]+}", cfg.wsHandler)

	log.Println("listening on http://localhost:8080")
	http.ListenAndServe("localhost:8080", r)
}

func initalizeCfg() (Config, error) {
	homeTemplate, err := template.ParseFiles("./templates/home.html")
	if err != nil {
		return Config{}, errors.New("error parsing home.html : " + err.Error())
	}
	loginTemplate, err := template.ParseFiles("./templates/login.html")
	if err != nil {
		return Config{}, errors.New("error parsing login.html : " + err.Error())
	}
	signupTemplate, err := template.ParseFiles("./templates/signup.html")
	if err != nil {
		return Config{}, errors.New("error parsing signup.html : " + err.Error())
	}
	chatTemplate, err := template.ParseFiles("./templates/chat.html")
	if err != nil {
		return Config{}, errors.New("error parsing chat.html : " + err.Error())
	}
	tmpl := templateConfig{
		home:   homeTemplate,
		login:  loginTemplate,
		signup: signupTemplate,
		chat:   chatTemplate,
	}
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	if err = godotenv.Load(); err != nil {
		return Config{}, errors.New("error .env file not found : " + err.Error())
	}
	JWT_SECRET := os.Getenv("JWT_SECRET")
	if JWT_SECRET == "" {
		return Config{}, errors.New("error JWT_SECRET not found in .env : " + err.Error())
	}
	MONGO_URI := os.Getenv("MONGO_URI")
	if MONGO_URI == "" {
		return Config{}, errors.New("error MONGO_URI not found in .env : " + err.Error())
	}

	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		return Config{}, errors.New("error could not connect to mongoDB : " + err.Error())
	}

	dbCfg := databaseConfig{
		mongoClient: mongoClient,
		userColl:    mongoClient.Database("sckt").Collection("user"),
		roomColl:    mongoClient.Database("sckt").Collection("room"),
		messageColl: mongoClient.Database("sckt").Collection("message"),
	}

	cfg := Config{
		db:         &dbCfg,
		tmpl:       &tmpl,
		upgrader:   upgrader,
		JWT_SECRET: []byte(JWT_SECRET),
		Clients:    []Client{},
	}
	return cfg, nil
}
