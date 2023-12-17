package main

import (
	"context"
	"log"
	"time"
)

type Room struct{
	Code int `bson:"code"`
	Messages []Message `bson:"messages"`  
	Users []User `bson:"users"` 
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}


func createRoom(code int){
	room := Room{
		Code: code,
		Messages: []Message{},
		Users: []User{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_,err := dbCfg.roomColl.InsertOne(context.TODO(),room)
	if err != nil {
		log.Println("error creating room : ",err)
	}
}

func getAllRooms(){

}

func deleteAllRooms(){

}