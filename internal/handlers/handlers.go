package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"passman/internal/auth"
	"passman/internal/db"
	"passman/internal/utils"
)

type req struct {
	Uname        string `json:"uname"`
	MasterPass   string `json:"master_pass"`
	SessionToken string `json:"session_token,omitempty"`
}

func HandleSignup(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var Req req

		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&Req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			log.Print("bad request: ", err)
			return
		}
		uname := Req.Uname
		pass := Req.MasterPass
		pass_hash, err := utils.HashPassword(pass)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Print("Internal server error in handlers.go :24")
			return
		}

		fmt.Printf("uname: %v\n", uname)
		fmt.Printf("master_pass: %v\n", pass)
		session_token, err := auth.GenerateSessionId()
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Print("Internal server error in handlers.go :43 : ", err)
			return
		}
		userParams := db.CreateUserParams{
			Username:           uname,
			MasterPasswordHash: pass_hash,
			SessionToken:       sql.NullString{String: session_token, Valid: true},
		}

		err = queries.CreateUser(context.Background(), userParams)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Print("Internal server error in handlers.go :36 : ", err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    session_token,
			MaxAge:   300,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "username",
			Value:    uname,
			MaxAge:   300,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
		})

		w.WriteHeader(http.StatusOK)
	}
}

func HandleSignin(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var Req req

		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&Req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			log.Print("bad request: ", err)
			return
		}
		uname := Req.Uname
		pass := Req.MasterPass

		user, err := queries.GetUserByUsername(context.Background(), uname)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Print("Internal server error in handlers.go :102 : ", err)
			return
		}

		isValidPass, err := utils.VerifyPassword(pass, user.MasterPasswordHash)
		if err != nil || !isValidPass {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			log.Print("Unauthorized access attempt in handlers.go :109 : ", err)
			return
		}

		sessionid, err := auth.GenerateSessionId()
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Print("Internal server error in handlers.go :116 : ", err)
			return
		}
		err = queries.InsertSessionTokenWithUsername(context.Background(), db.InsertSessionTokenWithUsernameParams{
			SessionToken: sql.NullString{String: sessionid, Valid: true},
			Username:     uname,
		})
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Print("Internal server error in handlers.go :125 : ", err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionid,
			MaxAge:   300,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
		})
		http.SetCookie(w, &http.Cookie{
			Name:     "username",
			Value:    uname,
			MaxAge:   300,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
		})
		w.WriteHeader(http.StatusOK)
	}
}
