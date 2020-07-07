package store

import (
	"context"
	"fmt"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/option"

	log "github.com/sirupsen/logrus"
)

// LikeDoc represents single document for a like
type LikeDoc struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Count   int32  `json:"count"`
	LikeMap map[string]bool
}

type FsPostDoc struct {
	Name    string `json:"name"`
	Count   int32  `json:"count"`
	LikeMap map[string]bool
}

type FireStoreStorage struct {
	client *firestore.Client
	ctx    context.Context
}

func (fs *FireStoreStorage) GetPostsMap() map[string]bool {

	posts, err := fs.GetAllPosts()

	if err != nil {
		log.Fatal("[GetPostsMap] Could not get all posts")
	}

	postMap := make(map[string]bool)
	post := &LikeDoc{}

	for _, v := range posts {
		v.DataTo(post)
		postMap[post.Name] = true
	}

	return postMap
}

func (fs *FireStoreStorage) ProcessPosts(name []string) {

	postMap := fs.GetPostsMap()

	log.Info("Processing %s posts", len(name))

	for _, v := range name {

		if _, ok := postMap[v]; ok {
			log.Print("%s exists! \n", v)
		} else {
			log.Print("%s does not exist", v)

			success, err := fs.AddNewPost(&LikeDoc{
				Name:    v,
				Count:   0,
				LikeMap: make(map[string]bool),
			})

			if err != nil || success != true {
				log.Print("Error writing doc")
				log.Print(err)
			}
		}
	}
}

func NewFireStoreStorage(ctx context.Context, projectName string, credLoc string) *FireStoreStorage {
	client, err := firestore.NewClient(ctx, projectName, option.WithCredentialsFile(credLoc))

	if err != nil {
		fmt.Println(err)
		log.Fatal("Could not get firestore client")
	}

	return &FireStoreStorage{
		client: client,
		ctx:    ctx,
	}
}

func (fs *FireStoreStorage) GetAllPosts() ([]*firestore.DocumentSnapshot, error) {
	posts, err := fs.client.Collection("posts").Documents(fs.ctx).GetAll()

	return posts, err
}

func (fs *FireStoreStorage) GetPostByName(name string) (*LikeDoc, error) {
	var doc *firestore.DocumentSnapshot
	var postDoc *FsPostDoc
	var err error

	p := fs.client.Collection("posts").
		Where("Name", "==", name).
		Documents(fs.ctx)

	doc, err = p.Next()

	if err != nil {
		return nil, err
	}

	doc.DataTo(&postDoc)
	postID := doc.Ref.ID

	likedoc := &LikeDoc{
		ID:      postID,
		Name:    postDoc.Name,
		Count:   postDoc.Count,
		LikeMap: postDoc.LikeMap,
	}

	return likedoc, nil
}

func (fs *FireStoreStorage) AddNewPost(doc *LikeDoc) (bool, error) {
	c := fs.client.Collection("posts")

	_, _, err := c.Add(fs.ctx, doc)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (fs *FireStoreStorage) LikePostByID(userId string, id string) (bool, error) {

	var err error
	dc, err := fs.client.Collection("posts").Doc(id).Get(fs.ctx)

	var p LikeDoc
	dc.DataTo(&p)

	p.LikeMap[userId] = true

	i := p.Count + 1

	log.Infof("[LikePostByID] Like for user %s | post %s | likes %d", userId, id, i)

	_, err = dc.Ref.Update(fs.ctx, []firestore.Update{
		{Path: "Count", Value: i},
		{Path: "LikeMap", Value: p.LikeMap},
	})

	if err != nil {
		return false, err
	}

	return true, nil
}
