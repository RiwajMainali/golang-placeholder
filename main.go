package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strings"
)

type dbData struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Passowrd string `json:"password"`
	Dbname   string `json:"dbname"`
}

var dbDatas dbData

type LoginDetails struct {
	Email string `json:"username"`
	//password string `json:password`
}
type LoginResult struct {
	UserID   int    `json:"UserID"`
	UserName string `json:"UserName"`
}

// returns the user details from DB as json
// func getUserID(data LoginDetails, db *sql.DB, w http.ResponseWriter) LoginResult, *string {
// 	rows, err2 := db.Query("select user_id, name from users where email=$1", data.Email)
// 	if err2 != nil {
// 		SendErrorMessage(err2.Error(), w)
// 		return
// 	}
// 	defer rows.Close()
// 	var user LoginResult
//
// 	for rows.Next() {
// 		err := rows.Scan(&user.UserID, &user.UserName)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// 	return user, nil
// }

// parses incoming data, sends to getUserId for returned data
func handleLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != "POST" {
		message := Response{
			Message: StrPtr("Please use POST"),
			Success: false,
			Data:    nil,
		}
		fmt.Fprintf(w, toJSON(message))
		return
	}
	var loginData LoginDetails
	err := json.NewDecoder(r.Body).Decode(&loginData)
	if err != nil {
		SendErrorMessage(err.Error(), w)
	}
	var ID int
	err3 := db.QueryRow("select user_id from users where email=$1", loginData.Email).Scan(&ID)
	if err3!=nil{
	SendErrorMessage(err3.Error(), w)
	}
	userToken, tokenErr := CreateToken(ID, db)
	if tokenErr != nil {
		SendErrorMessage(*tokenErr, w)
		return

	}
	response := Response{
		Message: nil,
		Success: true,
		Data: map[string]interface{}{
			"token": userToken.user_token,
			"time":  userToken.token_time,
		},
	}
	w.Header().Set("Content-type", "application/json")
	fmt.Fprintf(w, toJSON(response))
	return
}
func handleLogout(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendErrorMessage("Please use POST", w)
	}
}
func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		SendErrorMessage("Please use POST", w)
	}

}

type MyHandler struct {
	db *sql.DB
}

func (h MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch strings.ToLower(r.URL.Path) {
	case "/":
		fmt.Fprintln(w, "Welcome to home page")
	case "/login":
		handleLogin(w, r, h.db)
	case "/register":
		fmt.Print("/register")
		handleRegister(w, r)
	case "/logout":
		fmt.Print("/logout")
		handleLogout(h.db, w, r)
	default:
		http.NotFound(w, r)
	}
}

var GlobalDbDataError bool = false

func readData() {
	file, err := os.ReadFile("data.json")
	if err != nil {
		log.Fatal("cannot read data")
		return
	}
	err = json.Unmarshal(file, &dbDatas)
	if err != nil {
		GlobalDbDataError = true
	}

}
func main() {
	readData()
	psqlinfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbDatas.Host, 5432, dbDatas.User, dbDatas.Passowrd, dbDatas.Dbname)
	fmt.Printf(psqlinfo)
	db, err := sql.Open("postgres", psqlinfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	fmt.Println("Listen and served")
	h := MyHandler{
		db: db,
	}
	http.ListenAndServe(":8080", h)
}
