package main

import (
	"fmt"
	"log"
	"musicapp-server/db"
	"musicapp-server/middlewares"
	"musicapp-server/routes"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	db, err := db.ConnectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	/* 	_, err = db.Exec("DROP TABLE IF EXISTS users")
	   	if err != nil {
	   		log.Fatal(err)
	   	} else {
	   		log.Println("Table users dropped successfully")
	   	} */

	// Users Table
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS users (id UUID PRIMARY KEY ,name TEXT, email TEXT, password TEXT)")
	if err != nil {
		log.Fatal(err)
	}

	// Songs Table
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS songs(id SERIAL PRIMARY KEY,title VARCHAR(255) NOT NULL,artist VARCHAR(255) NOT NULL,album VARCHAR(255),genre VARCHAR(100),release_date DATE,duration INTEGER)")
	if err != nil {
		log.Fatal(err)
	}
	r := mux.NewRouter()

	// Handle Routes
	r.HandleFunc("/signup", routes.SignupHandler(db)).Methods("POST")
	r.HandleFunc("/login", routes.LoginHandler(db)).Methods("POST")
	r.Handle("/home", middlewares.JWTAuthMiddleware(http.HandlerFunc(routes.HomeHandler(db)))).Methods("GET")

	// Start Server
	fmt.Println("Server started on port :8000")
	log.Fatal(http.ListenAndServe(":8000", middlewares.JsonContentType(r)))
}
