package auth

import (
	"log"
	"net/http"
)

// Auth represents authentication middleware
type Auth struct {
	// Key used for public facing routes
	PublicAPIKey string
	// PrivateApiKey key used for private routes
	PrivateAPIKey string
}

// NewAuth creates a new instance of Auth
func NewAuth(publicAPIKey string, privateAPIKey string) *Auth {
	return &Auth{
		PublicAPIKey:  publicAPIKey,
		PrivateAPIKey: privateAPIKey,
	}
}

// PublicAuth authenticates using set public key
func (a *Auth) PublicAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Header.Get("PUBLIC_API_KEY") != a.PublicAPIKey {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("403 HTTP status code returned!"))
		} else {
			log.Println(r.RequestURI)
			next.ServeHTTP(w, r)
		}
	})
}

// PrivateAuth authenticates using set private key
func (a *Auth) PrivateAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("PRIVATE")
		if r.Header.Get("PRIVATE_API_KEY") != a.PrivateAPIKey {
			w.WriteHeader(http.StatusForbidden)
			w.Write([]byte("403 HTTP status code returned!"))
		} else {
			log.Println(r.RequestURI)
			next.ServeHTTP(w, r)
		}
	})
}
