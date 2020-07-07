package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	actions "github.com/olingern/blog-backend/pkg/actions"
	"github.com/olingern/blog-backend/pkg/auth"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := mux.NewRouter()

	amw := auth.NewAuth(os.Getenv("PUBLIC_API_KEY"), os.Getenv("PRIVATE_API_KEY"))

	ctx := context.Background()
	actions := actions.NewActions(ctx)

	webhookHandler := http.HandlerFunc(actions.Scrape)

	r.HandleFunc("/users", actions.CreateID).Methods("POST", "OPTIONS")
	r.HandleFunc("/users/{userId}/posts/{name}", actions.GetPostInfo).Methods("GET", "OPTIONS")
	r.HandleFunc("/users/{userId}/likes/{id}", actions.Like).Methods("POST", "OPTIONS")

	r.Handle("/webhooks/netlify", amw.PrivateAuth(webhookHandler)).Methods("POST")

	http.ListenAndServe(":8080", handlers.RecoveryHandler()(r))

}
