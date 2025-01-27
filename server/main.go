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

	r := mux.NewRouter()

	// Handle Routes
	r.HandleFunc("/signup", routes.SignupHandler(db)).Methods("POST")
	r.HandleFunc("/login", routes.LoginHandler(db)).Methods("POST")
	r.Handle("/home", middlewares.JWTAuthMiddleware(http.HandlerFunc(routes.HomeHandler(db)))).Methods("GET")
	r.Handle("/home/songs", routes.SongHandler(db)).Methods("GET")

	// Start Server
	fmt.Println("Server started on port :8000")
	log.Fatal(http.ListenAndServe(":8000", middlewares.JsonContentType(r)))
}
