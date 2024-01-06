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

type Room struct {
	Name      string    `bson:"name"`
	Code      int       `bson:"code"`
	Messages  []Message `bson:"messages"`
	Users     []User    `bson:"users"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (db *databaseConfig) generateRandomRoomCode() int {
	var num int
	for {
		num = rand.Intn(10000) + 100000
		if _, err := db.getRoomByCode(num); err != nil {
			break
		}
	}
	return num
}

func (db *databaseConfig) getRoomByCode(code int) (Room, error) {
	var room Room
	result := db.roomColl.FindOne(context.TODO(), bson.D{{Key: "code", Value: code}}).Decode(room)
	if result == mongo.ErrNoDocuments {
		return Room{}, errors.New("room not found")
	}
	return room, nil
}

func (db *databaseConfig) createRoom(user User, name string) (Room, error) {
	room := Room{
		Name:      name,
		Code:      db.generateRandomRoomCode(),
		Messages:  []Message{},
		Users:     []User{user},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := db.roomColl.InsertOne(context.TODO(), room)
	if err != nil {
		log.Println("error creating room : ", err)
		return Room{}, err
	}
	return room, nil
}

func (db *databaseConfig) deleteAllRooms() {
	result, err := db.roomColl.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		log.Println("error in deleting rooms: ", err)
	}
	log.Println("Documents delete : ", result.DeletedCount)
}

func (db *databaseConfig) getAllRooms() {
	cursor, err := db.roomColl.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println("error in finding rooms : ", err.Error())
	}
	var rooms []bson.M
	if err := cursor.All(context.TODO(), &rooms); err != nil {
		log.Println("error quering rooms : ", err.Error())
	}

	for _, room := range rooms {
		log.Println(room)
	}
}
