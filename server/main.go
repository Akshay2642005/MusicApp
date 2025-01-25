package main

import (
	"fmt"
	"log"
	"musicapp-server/db"
	"musicapp-server/routes"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	pgdb, err := db.ConnectToDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer pgdb.Close()

	_, err = pgdb.Exec("CREATE TABLE IF NOT EXISTS users (name TEXT, email TEXT, password TEXT)")
	if err != nil {
		log.Fatal(err)
	}
	r := mux.NewRouter()

	// Handle Routes

	r.HandleFunc("/signup", routes.SignupHandler(pgdb)).Methods("POST")
	r.HandleFunc("/login", routes.LoginHandler(pgdb)).Methods("POST")

	// Start Server
	fmt.Println("Server started on port :8000")
	log.Fatal(http.ListenAndServe(":8000", jsonContentTypeMiddleware(r)))
}

func jsonContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}