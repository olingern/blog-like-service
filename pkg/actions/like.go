package actions

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"

	log "github.com/sirupsen/logrus"
)

// LikeAPIResponse represents a like API response
type LikeAPIResponse struct {
	DocID string `json:"docId"`
	ID    string `json:"id"`
	Count int32  `json:"count"`
}

type GetPostIDResponse struct {
	PostId   string `json:"postId"`
	Likes    int32  `json:"likes"`
	HasLiked bool   `json:"hasLiked"`
}

func write5xx(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 - Internal Server Error"))
}

func write4xx(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("400 - Bad Request"))
}

func createIDFromReq(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")

	if forwarded != "" {
		return forwarded
	}

	return r.RemoteAddr + "-" + r.Header.Get("User-Agent")
}

type CreateIDResponse struct {
	UserID string `json:"userId"`
}

// handles POST /users
func (a *Actions) CreateID(w http.ResponseWriter, r *http.Request) {
	// TODO: move to common middleware. fix security issue of origin allow all.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "public_api_key")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	res := &CreateIDResponse{
		UserID: uuid.New().String(),
	}

	json.NewEncoder(w).Encode(res)
}

// GetPostId handles GET /:name
func (a *Actions) GetPostInfo(w http.ResponseWriter, r *http.Request) {
	// TODO: move to common middleware. fix security issue of origin allow all.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Headers", "public_api_key")

	vars := mux.Vars(r)

	name := ConvertName(vars["name"])
	log.Infof("[GetPostInfo]: name: %s", name)

	likeDoc, err := a.store.GetPostByName(name)

	if err != nil {
		log.Errorf("[GetPostInfo]: a.store.GetPostByName error: %s", err)
		write5xx(w)
		return
	}

	_, ok := likeDoc.LikeMap[vars["userId"]]

	res := &GetPostIDResponse{
		PostId:   likeDoc.ID,
		Likes:    likeDoc.Count,
		HasLiked: ok,
	}

	json.NewEncoder(w).Encode(res)
}

// Like handles POST /:id/like
func (a *Actions) Like(w http.ResponseWriter, r *http.Request) {

	// TODO: move to common middleware. fix security issue of origin allow all.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "public_api_key")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	log.Infof("[Like] METHOD: %s", r.Method)

	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}

	vars := mux.Vars(r)

	if vars["userId"] == "" {
		w.WriteHeader(403)
	}

	success, err := a.store.LikePostByID(vars["userId"], vars["id"])

	if err != nil {
		log.Print(err)
		write5xx(w)
		return
	}

	if success != true {
		log.Printf("Could not upvote post %s for user %s", vars["id"], vars["userId"])
		write5xx(w)
		return
	}

	w.WriteHeader(200)
}
