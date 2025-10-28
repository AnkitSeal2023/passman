package main

import (
	"log"
	"net/http"
	"passman/internal/auth"
	"passman/internal/handlers"
	"passman/views/pages"
)

func serveEntities(w http.ResponseWriter, r *http.Request) {
	if err := pages.EntitiesListPage().Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func serveSubEntities(w http.ResponseWriter, r *http.Request) {
	if err := pages.SubEntities().Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func serveSignupPage(w http.ResponseWriter, r *http.Request) {
	if err := pages.SignupPage().Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func serveSigninPage(w http.ResponseWriter, r *http.Request) {
	if err := pages.SigninPage().Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("views/static"))))
	mux.HandleFunc("/", auth.RequreAuth(serveEntities))
	mux.HandleFunc("/subentities", auth.RequreAuth(serveSubEntities))
	mux.HandleFunc("/signup", serveSignupPage)
	mux.HandleFunc("/signin", serveSigninPage)

	mux.HandleFunc("/api/signup", handlers.HandleSignup)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
