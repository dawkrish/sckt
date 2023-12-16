package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func conncectToDB() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
	MONGO_URI := os.Getenv("MONGO_URI")
	if MONGO_URI == "" {
		log.Fatal("MONGO_URI not set")
	}

	var err error
	dbCfg.mongoClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(MONGO_URI))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("mongo connection established")

	dbCfg.userColl = dbCfg.mongoClient.Database("sckt").Collection("user")
}
