package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"passman/internal/auth"
	"passman/internal/db"
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
		log.Fatal("serving signup page err")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func serveSigninPage(w http.ResponseWriter, r *http.Request) {
	if err := pages.SigninPage().Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	connStr := "postgresql://postgres:2004@localhost:5432/passman?sslmode=disable"
	conn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	queries := db.New(conn)

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("views/static"))))
	mux.HandleFunc("/", auth.RequreAuth(serveEntities, queries))
	mux.HandleFunc("/subentities", auth.RequreAuth(serveSubEntities, queries))
	mux.HandleFunc("/signup", serveSignupPage)
	mux.HandleFunc("/signin", serveSigninPage)
	mux.HandleFunc("/api/signup/new", handlers.HandleSignup(queries))
	mux.HandleFunc("/api/signin", handlers.HandleSignin(queries))

	log.Fatal(http.ListenAndServe(":5000", mux))
}
