package main

import (
	"fmt"
	"net/http"

	"github.com/Shanmukh1609/backend/models"
	"github.com/Shanmukh1609/backend/routers"
)

func main(){
	fmt.Println("Welcome!!");
	models.InitDB()
	fmt.Println("Database is connected successfully");
	fmt.Println("Starting the Game")
	mux:=http.NewServeMux();

	routers.GameEntry(mux);
	http.ListenAndServe(":8080", mux);

}

