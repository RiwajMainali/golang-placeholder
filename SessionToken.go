package main

import (
	"database/sql"
	"fmt"
	"net/http"
)

func CreateToken(userID int, db *sql.DB, w http.ResponseWriter, r *http.Request){
	rows, err := db.Query("select user_id from users where user_token is null and user_id=$1;", userID)
	if err!=nil{
		w.Header().Set("Content-Type","application/json")
		fmt.Fprintf(w, `{"message":"User already Authenticated", "success":"True"}`)
		return
	}
	
}

