package main

import "time"

type Message struct {
	Text      string `bson:"text"`
	RoomCode  int    `bson:"room_code"`
	Sender    User   `bson:"sender"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}


func createMessage(){

}

func deleteMessage(){

}

func updateMessage(){

}

func getAllMessages(){

}

func deleteAllMessages(){
	
}