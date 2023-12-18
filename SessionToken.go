package main

import (
	"database/sql"
	"fmt"
	_ "math/big"
	"net/http"

	_ "github.com/google/uuid"
)

type token struct {
	user_token *string
	token_time *string 
}

func strPtr(s string) *string {
    return &s
}
func GetToken(userID int, db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		message := Response{
			Message: strPtr("Please use GET"),
			Success: false,
			Data:    nil,
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, toJSON(message))
		return
	}
	rows, err := db.Query("select user_token, token_time from users where user_token is not null and user_id=$1;", userID)
	if err != nil {
		message := Response{
			Message: strPtr("User already Authenticated"),
			Success: false,
			Data:    nil,
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, toJSON(message))
		return
	}
	var tokenData token
	for rows.Next() {
		err := rows.Scan(&tokenData.user_token, &tokenData.token_time)
		if err != nil {
			message := Response{
				Message: strPtr(err.Error()),
				Success: false,
				Data:    nil,
			}

			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, toJSON(message))
			return
		}
		message := Response{
			Message: nil,
			Success: true,
			Data: map[string]interface{}{
				"expirationToken": tokenData.token_time,
				"userToken":       tokenData.user_token,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, toJSON(message))
		return
	}
}
