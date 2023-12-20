package main

import (
	"database/sql"
	"fmt"
	_ "math/big"
	"net/http"
	"strconv"
	"github.com/google/uuid"
	"time"
)

type token struct {
	user_token *string
	token_time *string
}

func ClearToken(userID int, db *sql.DB) bool {
	_, err := db.Exec("update users set token_time= null, user_token=null;")
	if err != nil {
		return false
	}
	return true
}
func CheckStatus(userID int, db *sql.DB)(*bool, *string){
	var exists bool
	err := db.QueryRow("SELECT EXISTS (SELECT * FROM users WHERE user_token is null and user_id=$1);", userID).Scan(&exists)
	if err != nil {
		return &exists, StrPtr("Unable to create session")
	}
	return &exists, nil

}
func CreateToken(userID int, db *sql.DB) (*token, *string) {
	tokenTime := time.Now().Add(24 * time.Hour).Unix()
	userToken := uuid.New()
	_, err := db.Exec("update users set token_time=$1, user_token=$2 where user_id=$3;", tokenTime, userToken, userID)
	if err != nil {
		return nil, StrPtr("Unable to create session")
	}
	tokenRes :=token{
		user_token: StrPtr(userToken.String()),
		token_time: StrPtr(strconv.FormatInt(tokenTime,10)) ,
	}
	return &tokenRes, nil
}
func GetToken(userID int, db *sql.DB) (*token, *string) {
	rows, err := db.Query("select user_token, token_time from users where user_token is not null and user_id=$1;", userID)
	defer rows.Close()
	if err != nil {
		return nil, StrPtr("DB error")
	}
	var tokenData token
	for rows.Next() {
		err := rows.Scan(&tokenData.user_token, &tokenData.token_time)
		if err != nil {
			return nil, StrPtr("DB scan error")
		}
		return &tokenData, nil
	}
	return nil, StrPtr("Massive error: User not found")
}
func SendToken(userID int, db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		message := Response{
			Message: StrPtr("Please use GET"),
			Success: false,
			Data:    nil,
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, toJSON(message))
		return
	}
	
	rows, err := db.Query("select user_token, token_time from users where user_token is not null and user_id=$1;", userID)
	defer rows.Close()
	if err != nil {
		message := Response{
			Message: StrPtr("User already Authenticated"),
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
				Message: StrPtr(err.Error()),
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
