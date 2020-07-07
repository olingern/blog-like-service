package actions

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"

	"cloud.google.com/go/firestore"
	store "github.com/olingern/blog-backend/pkg/store"
)

// Actions is externally available obj
type Actions struct {
	collectionName string
	client         *firestore.Client
	ctx            context.Context
	projectName    string
	siteURL        string
	store          *store.FireStoreStorage
}

// NewActions creates instance of Actions
func NewActions(ctx context.Context) *Actions {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	projectName := os.Getenv("GCP_PROJECT_NAME")
	credLoc := os.Getenv("GCP_CREDENTIAL_LOC")
	siteURL := os.Getenv("SITE_URL")

	store := store.NewFireStoreStorage(ctx, projectName, credLoc)

	store.GetAllPosts()
	store.GetPostByName("test")

	return &Actions{
		store:          store,
		collectionName: "posts",
		ctx:            ctx,
		projectName:    projectName,
		siteURL:        siteURL,
	}
}
