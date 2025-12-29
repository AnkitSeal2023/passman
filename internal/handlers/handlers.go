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
	"passman/views/pages"
)

type req struct {
	Uname      string `json:"uname"`
	MasterPass string `json:"master_pass"`
}
type newCredentialReq struct {
	Website      string `json:"website"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	PasswordName string `json:"passwordName,omitempty"`
}

func HandleSignup(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

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

		dek, err := utils.GenerateDEK()
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Print("Internal server error in handlers.go :83 : ", err)
			return
		}
		encryptedDEK, err := utils.EncryptDEKWithKEK(pass, dek)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Print("Internal server error in handlers.go :89 : ", err)
			return
		}

		userParams := db.CreateUserParams{
			Username:           uname,
			MasterPasswordHash: pass_hash,
			SessionToken:       sql.NullString{String: session_token, Valid: true},
			EncrDek:            encryptedDEK,
		}

		err = queries.CreateUser(context.Background(), userParams)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Print("Internal server error in handlers.go :65 : ", err)

			return
		}

		// Get the newly created user to get their userid
		user, err := queries.GetUserByUsername(context.Background(), uname)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Print("Error getting user after creation: ", err)
			return
		}

		// Create a default vault for the user
		maxVaultIdRaw, err := queries.GetMaxVaultId(context.Background())
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Print("Error getting max vault id: ", err)
			return
		}

		maxVaultId, ok := maxVaultIdRaw.(int64)
		if !ok {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Print("Error converting max vault id to int64")
			return
		}

		newVaultId := int32(maxVaultId + 1)
		vaultName := uname + "'s vault"

		err = queries.CreateVault(context.Background(), db.CreateVaultParams{
			VaultID:                newVaultId,
			Userid:                 user.Userid,
			VaultItemNameEncrypted: vaultName,
		})
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Print("Error creating vault: ", err)
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
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

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

func HandleNewCredential(queries *db.Queries) http.HandlerFunc {
	passPhrase := "correct horse battery staple"

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		cookie, err := r.Cookie("username")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			log.Print("No username cookie: ", err)
			return
		}
		username := cookie.Value

		var Req newCredentialReq
		defer r.Body.Close()
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&Req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			log.Print("bad request: ", err)
			return
		}
		website := Req.Website
		if website == "" && Req.PasswordName != "" {
			website = Req.PasswordName
		}
		usernameForSite := Req.Username
		password := Req.Password

		encryptedWebsiteName, err := utils.EncryptUsingPassphrase(passPhrase, []byte(website))
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Print("Error encrypting website name: ", err)
			return
		}

		encryptedUsername, err := utils.EncryptUsingPassphrase(passPhrase, []byte(usernameForSite))
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Print("Error encrypting username: ", err)
			return
		}

		encryptedPassword, err := utils.EncryptUsingPassphrase(passPhrase, []byte(password))
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Print("Error encrypting password: ", err)
			return
		}

		user, err := queries.GetUserByUsername(context.Background(), username)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Print("Error getting user: ", err)
			return
		}

		vault, err := queries.GetVaultByUserId(context.Background(), user.Userid)
		if err != nil {
			log.Print("Vault doesn't exist for user, creating one: ", err)

			maxVaultIdRaw, err := queries.GetMaxVaultId(context.Background())
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				log.Print("Error getting max vault id: ", err)
				return
			}

			maxVaultId, ok := maxVaultIdRaw.(int64)
			if !ok {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				log.Print("Error converting max vault id to int64")
				return
			}

			newVaultId := int32(maxVaultId + 1)
			vaultName := username + "'s vault"

			err = queries.CreateVault(context.Background(), db.CreateVaultParams{
				VaultID:                newVaultId,
				Userid:                 user.Userid,
				VaultItemNameEncrypted: vaultName,
			})
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				log.Print("Error creating vault: ", err)
				return
			}

			vault, err = queries.GetVaultByUserId(context.Background(), user.Userid)
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				log.Print("Error getting newly created vault: ", err)
				return
			}
		}

		var vaultItem db.GetVaultItemsForUserRow
		found := false
		userVaultItems, err := queries.GetVaultItemsForUser(context.Background(), username)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Print("Error listing vault items for user: ", err)
			return
		}
		for _, vi := range userVaultItems {
			nameBytes, err := utils.DecryptUsingPassphrase(passPhrase, vi.VaultItemNameEncrypted)
			if err != nil {
				continue
			}
			if string(nameBytes) == website {
				vaultItem = vi
				found = true
				break
			}
		}
		if !found {
			maxVaultItemIdRaw, err := queries.GetMaxVaultItemId(context.Background())
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				log.Print("Error getting max vault item id: ", err)
				return
			}

			maxVaultItemId, ok := maxVaultItemIdRaw.(int64)
			if !ok {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				log.Print("Error converting max vault item id to int64")
				return
			}

			newVaultItemId := int32(maxVaultItemId + 1)
			err = queries.CreateVaultItem(context.Background(), db.CreateVaultItemParams{
				VaultItemID:            newVaultItemId,
				VaultID:                vault.VaultID,
				VaultItemNameEncrypted: encryptedWebsiteName,
			})
			if err != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				log.Print("Error creating vault item: ", err)
				return
			}
			vaultItem.VaultItemID = newVaultItemId
			vaultItem.VaultItemNameEncrypted = encryptedWebsiteName
		}

		maxCredentialIdRaw, err := queries.GetMaxCredentialId(context.Background())
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Print("Error getting max credential id: ", err)
			return
		}

		maxCredentialId, ok := maxCredentialIdRaw.(int64)
		if !ok {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Print("Error converting max credential id to int64")
			return
		}

		newCredentialId := int32(maxCredentialId + 1)
		err = queries.CreateCredential(context.Background(), db.CreateCredentialParams{
			CredentialID:          newCredentialId,
			VaultItemID:           vaultItem.VaultItemID,
			Encryptedcredname:     encryptedUsername,
			Encryptedcredpassword: encryptedPassword,
		})
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Print("Error creating credential: ", err)
			return
		}

		log.Printf("Successfully created credential for user %s", username)
		w.WriteHeader(http.StatusOK)
	}
}

func HandleGetCredentialsByEntity(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			w.Header().Set("Allow", http.MethodGet)
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		entityName := r.URL.Query().Get("account")
		if entityName == "" {
			http.Error(w, "entity name is required", http.StatusBadRequest)
			return
		}

		cookie, err := r.Cookie("username")
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			log.Print("No username cookie: ", err)
			return
		}
		username := cookie.Value

		userDEK, err := queries.GetUserDEK(context.Background(), username)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Print("Error getting user DEK: ", err)
			return
		}

		if userDEK == "" {
			http.Error(w, "user DEK not found", http.StatusInternalServerError)
			log.Print("User DEK not valid")
			return
		}

		vaultItems, err := queries.GetVaultItemsForUser(context.Background(), username)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			log.Print("Error fetching vault items: ", err)
			return
		}

		var allCredentials []pages.Credentials
		for _, vaultItem := range vaultItems {
			nameBytes, err := utils.DecryptUsingPassphrase("correct horse battery staple", vaultItem.VaultItemNameEncrypted)
			if err != nil || string(nameBytes) != entityName {
				continue
			}
			credentials, err := queries.GetAllCredentials(context.Background(), vaultItem.VaultItemID)
			if err != nil {
				log.Print("Error fetching credentials for vault item: ", err)
				continue
			}

			for _, cred := range credentials {
				decryptedName, err := utils.DecryptUsingPassphrase("correct horse battery staple", cred.Encryptedcredname)
				if err != nil {
					log.Print("Error decrypting credential name: ", err)
					continue
				}

				decryptedPass, err := utils.DecryptUsingPassphrase("correct horse battery staple", cred.Encryptedcredpassword)
				if err != nil {
					log.Print("Error decrypting password: ", err)
					continue
				}

				allCredentials = append(allCredentials, pages.Credentials{
					UserName: string(decryptedName),
					Pass:     string(decryptedPass),
				})
			}
		}

		component := pages.SubEntities(entityName, allCredentials)
		if err := component.Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
