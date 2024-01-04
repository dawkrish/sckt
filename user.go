package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	DisplayName string    `bson:"display_name"`
	UserName    string    `bson:"user_name"`
	Email       string    `bson:"email"`
	Password    string    `bson:"password"`
	CreatedAt   time.Time `bson:"created_at"`
	UpdatedAt   time.Time `bson:"updated_at"`
}

func (db *databaseConfig) createUser(userName, displayName, email, password string) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		log.Println("err in hashing password ", err)
	}
	user := User{
		UserName:    userName,
		DisplayName: displayName,
		Email:       email,
		Password:    string(hashedPassword),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	_, err = db.userColl.InsertOne(context.TODO(), user)
	if err != nil {
		log.Println("error while creating user: ", err)
		return User{}, err
	}
	return user, nil
}

func (db *databaseConfig) getUserByUserName(username string) (User, error) {
	filter := bson.D{{Key: "user_name", Value: username}}
	var user User
	result := db.userColl.FindOne(context.TODO(), filter).Decode(&user)
	if result == mongo.ErrNoDocuments {
		log.Println("user not found")
		return User{}, errors.New("user not found")
	}
	return user, nil

}

func (db *databaseConfig) getAllUsers() {
	cursor, err := db.userColl.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Println("error in finding users", err)
	}
	var users []bson.M
	if err := cursor.All(context.TODO(), &users); err != nil {
		log.Println("error quering users : ", err)
	}

	for _, user := range users {
		fmt.Println(user)
	}
}

func (db *databaseConfig) deleteAllUsers() {
	result, err := db.userColl.DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		log.Println("error in deleting users : ", err)
	}
	log.Println("Documents delete : ", result.DeletedCount)

}
