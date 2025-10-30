package handlers

import "net/http"

func HandleSignup(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	// uname := r.FormValue("uname")
	// pass := r.FormValue("pass")
}
