package main

import (
	"fmt"
	"net/http"
	"passman/views"
	"time"
)

func NewNowHandler(now func() time.Time) NowHandler {
	return NowHandler{Now: now}
}

type NowHandler struct {
	Now func() time.Time
}

func (nh NowHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	views.AccountsListPageComponent().Render(r.Context(), w)
}

func main() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("views"))))
	http.Handle("/", NewNowHandler(time.Now))

	fmt.Println("Starting server on port 8080")
	http.ListenAndServe(":8080", nil)
}
