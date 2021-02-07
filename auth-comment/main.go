package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // using postgres sql
)

var db *gorm.DB

func main() {
	initDB()
	http.HandleFunc("/register", Register)
	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/refresh", Refresh)
	http.HandleFunc("/createposts", CreatePost)
	http.HandleFunc("/addcomments", AddComments)
	http.HandleFunc("/deletecomments", DeleteComments)
	http.HandleFunc("/deleteposts", DeletePosts)
	//add sub comments / likes and reactions as well
	http.HandleFunc("/postinteraction", Interact)
	http.HandleFunc("/getinteraction", GetInteraction)
	//need to add alter comment and likes for the comment
	http.HandleFunc("/addreaction", AddReaction)
	http.HandleFunc("/getposts", GetPosts)
	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8001", nil))
}

func initDB() {
	var err error
	prosgretConname := fmt.Sprintf("dbname=%v password=%v port=%v user=%v host=%v sslmode=disable", os.Getenv("DB_NAME"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_HOST"))
	fmt.Printf("connecting to db on %v", os.Getenv("DB_HOST"))

	db, err = gorm.Open("postgres", prosgretConname)
	if err != nil {
		log.Printf("the error is :%v", err)
		panic("Failed to connect to database!")
	}
}
