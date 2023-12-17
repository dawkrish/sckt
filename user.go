package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type User struct{
	Name string `bson:"name"`
	Email string `bson:"email"`
	Password string `bson:"password"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func createUser(name, email , password string){
	user := User{
		Name : name,
		Email: email,
		Password: password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err := dbCfg.userColl.InsertOne(context.TODO(),user)
	if err != nil {
		log.Println("error while creating user: ",err)
	}
}

func getAllUsers(){
	cursor, err := dbCfg.userColl.Find(context.TODO(),bson.M{})
	if err != nil {
		log.Println("error in finding users",err)
	}
	var users []bson.M
	if err := cursor.All(context.TODO(),&users); err != nil {
		log.Println("error quering users : ",err)
	}
	log.Println(users[0]["name"])
}

func deleteAllUsers(){
	result,err := dbCfg.userColl.DeleteMany(context.TODO(),bson.M{})
	if err != nil {
		log.Println("error in deleting users : ",err)
	}
	log.Println("Documents delete : ",result.DeletedCount)

}

