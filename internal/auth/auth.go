package auth

import (
	"net/http"

	"passman/internal/db"
)

func RequreAuth(next http.HandlerFunc, queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		uname, err := r.Cookie("username")
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}

		err = isValidSessionToken(queries, c.Value, uname.Value)
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		next(w, r)
	}
}
