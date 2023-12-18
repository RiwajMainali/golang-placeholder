package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/lib/pq"
)

type dbData struct {
	Host     string `json:"host"`
	User     string `json:"user"`
	Passowrd string `json:"password"`
	Dbname   string `json:"dbname"`
}

func returnData(w http.ResponseWriter, r *http.Request){
	return 
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
func getUserID(data LoginDetails, db *sql.DB) LoginResult {
	rows, err2 := db.Query("select user_id, name from users where email=$1", data.Email)
	if err2 != nil {
		log.Fatal(err2)
	}
	defer rows.Close()
	var user LoginResult

	for rows.Next() {
		err := rows.Scan(&user.UserID, &user.UserName)
		if err != nil {
			log.Fatal(err)
		}
	}
	return user

}

// parses incoming data, sends to getUserId for returned data
func handleLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method == "POST" {
		var loginData LoginDetails
		err := json.NewDecoder(r.Body).Decode(&loginData)
		if err != nil {
			log.Fatal(err)
		}
		var user LoginResult
		user = getUserID(loginData, db)
		res, err := json.Marshal(user)
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Set("Content-type", "application/json")
		fmt.Fprintf(w, string(res))
	}
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		return
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
	case "/token":
		fmt.Print("/token")
		GetToken(1, h.db, w, r)
	default:
		http.NotFound(w, r)
	}
}

var GlobalDbDataError bool = false

func readData() {
	file, err := os.ReadFile("data.json")
	if err != nil {
		GlobalDbDataError = true
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
