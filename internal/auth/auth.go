package auth

import (
	"fmt"
	"net/http"
)

func RequreAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_id")
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}

		//logic to check if c is valid
		err = isValidSessionId(c.Name)
		if err != nil {
			http.Redirect(w, r, "/signin", http.StatusFound)
			return
		}
		fmt.Printf("%v\n", c)
		next(w, r)
	}
}
