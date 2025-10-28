package main

import (
	"log"
	"net/http"
	"passman/views/pages"
)

func ServeEntities(w http.ResponseWriter, r *http.Request) {
	if err := pages.EntitiesListPage().Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ServeSubEntities(w http.ResponseWriter, r *http.Request) {
	if err := pages.SubEntities().Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("views/static"))))
	http.HandleFunc("/", ServeEntities)
	http.HandleFunc("/subentities", ServeSubEntities)
	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
