package main

import "database/sql"

func CreateToken(userID int, db *sql.DB){
	rows, err := db.Query("select * from users where user_token is null and user_id=$1;", userID)
	if err!=nil{
		
	}
}

