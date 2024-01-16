package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type Message struct {
	Text      string    `bson:"text"`
	RoomCode  int       `bson:"room_code"`
	Sender    string    `bson:"sender"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (db *databaseConfig) createMessage(roomCode int, text string, sender string) (Message, error) {
	message := Message{
		RoomCode:  roomCode,
		Text:      text,
		Sender:    sender,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := db.messageColl.InsertOne(context.TODO(), message)
	if err != nil {
		log.Println("error creating message : ", err)
		return Message{}, err
	}
	return message, nil
}

func (db *databaseConfig) deleteMessage() {

}

func (db *databaseConfig) updateMessage() {

}

func (db *databaseConfig) getAllMessagesByRoomCode(roomCode int) ([]Message, error) {
	filter := bson.D{{Key: "room_code", Value: roomCode}}
	cursor, err := db.messageColl.Find(context.TODO(), filter)
	if err != nil {
		log.Println("error creating message cursor : ", err.Error())
		return nil, err
	}
	var messages []Message
	if err := cursor.All(context.TODO(), &messages); err != nil {
		log.Println("error quering rooms : ", err.Error())
		return []Message{}, err
	}

	return messages, nil

}

func (db *databaseConfig) deleteAllMessages() {

}
