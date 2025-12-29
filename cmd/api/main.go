package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	"passman/internal/auth"
	"passman/internal/db"
	"passman/internal/handlers"
	"passman/internal/utils"
	"passman/views/pages"
)

func serveEntities(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get username from cookie
		cookie, err := r.Cookie("username")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			log.Print("No username cookie: ", err)
			return
		}
		username := cookie.Value

		encEntities, err := queries.GetAllVaultItemsByUser(context.Background(), username)
		if err != nil {
			log.Print("Error fetching entities: ", err)
			encEntities = []string{}
		}

		var entities []string
		for _, enc := range encEntities {
			plain, err := utils.DecryptUsingPassphrase("correct horse battery staple", enc)
			if err != nil {
				continue
			}
			entities = append(entities, string(plain))
		}

		if err := pages.EntitiesListPage(entities).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
	mux.HandleFunc("/", auth.RequreAuth(serveEntities(queries), queries))
	mux.HandleFunc("/subentities", auth.RequreAuth(handlers.HandleGetCredentialsByEntity(queries), queries))
	mux.HandleFunc("/signup", serveSignupPage)
	mux.HandleFunc("/signin", serveSigninPage)
	mux.HandleFunc("/api/signup/new", handlers.HandleSignup(queries))
	mux.HandleFunc("/api/newcredential", handlers.HandleNewCredential(queries))
	mux.HandleFunc("/api/signin", handlers.HandleSignin(queries))

	log.Fatal(http.ListenAndServe(":5000", mux))
}
