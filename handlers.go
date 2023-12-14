package main

import (
	"log"
	"net/http"
	"html/template"
)

func wsHandler(w http.ResponseWriter, r *http.Request){
		log.Println("websocket connection established...")
		conn , err := upgrader.Upgrade(w,r,nil)
		if err != nil {
			log.Println("error in upgrading : ",err)
		}
		defer conn.Close()

		// infinite read loop
		for {
			mt, message, err := conn.ReadMessage()
			if err != nil {
				log.Println("read failed: ",err)
			}
			log.Println("message recieved : ",string(message))
			conn.WriteMessage(mt,message)
		}
}

func homeHandler(w http.ResponseWriter, r *http.Request){
		t,_ := template.ParseFiles("static/home.html")
		type dataToSend struct{
			Title string
		}
		d := dataToSend{
			Title: "Home Page",
		}
		t.Execute(w,d)
}