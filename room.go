package main

import (
	"context"
	"errors"
	"log"
	"math/rand"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Room struct{
	Name string `bson:"name"`
	Code int `bson:"code"`
	Messages []Message `bson:"messages"`  
	Users []User `bson:"users"` 
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (db *databaseConfig)generateRandomRoomCode() int {
    var num int
    for {
        num = rand.Intn(10)
        if _, err := db.getRoomByCode(num); err == nil {
            break
        }
    }
    return num
}

func (db *databaseConfig)getRoomByCode(code int) (Room, error){
	var room Room
	result := db.roomColl.FindOne(context.TODO(),bson.D{{Key: "code" , Value: code}}).Decode(room)
	if result == mongo.ErrNoDocuments{
		return Room{},errors.New("room not found")
	}
	return room, nil
}

func (db *databaseConfig)createRoom(name string){
	room := Room{
		Name : name,
		Code: db.generateRandomRoomCode(),
		Messages: []Message{},
		Users: []User{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_,err := db.roomColl.InsertOne(context.TODO(),room)
	if err != nil {
		log.Println("error creating room : ",err)
	}
}

func getAllRooms(){

}

func deleteAllRooms(){

}