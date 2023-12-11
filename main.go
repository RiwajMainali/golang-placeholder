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
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Passowrd string `json:"password"`
	Dbname string `json:"dbname"`
}


type LoginDetails struct{
Email string `json:"username"`
//password string `json:password`
}
type LoginResult struct {
	UserID int `json:"UserID"`
	UserName string `json:"UserName"`
}
func getUserID(data LoginDetails)LoginResult{
	psqlinfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		dbData.Host, dbData.Port, dbData.User, dbData.Passowrd, dbData.Dbname)
	db, err := sql.Open("postgres", psqlinfo)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err2 := db.Query("select user_id, name from users where email=$1", data.Email)
	if err2!=nil{
		log.Fatal(err2)
	}
	defer rows.Close()
	var user LoginResult

	for rows.Next(){

		err := rows.Scan(&user.UserID, &user.UserName)
		if(err!=nil){
			log.Fatal(err)
		}
	}
	return user

}

func handleLogin(w http.ResponseWriter, r *http.Request){
	if r.Method=="POST"{
		var loginData LoginDetails
		err :=json.NewDecoder(r.Body).Decode(&loginData)
		if (err!=nil){
			log.Fatal(err)
		}
		var user LoginResult
		user= getUserID(loginData)
		res, err := json.Marshal(user)
		if err!=nil{
			log.Fatal(err)
		}
		w.Header().Set("Content-type", "application/json")
		fmt.Fprintf(w,string(res))
	}
}

func handleRegister(w http.ResponseWriter, r *http.Request){
	if r.Method =="POST"{
		return
	}

}

type MyHandler struct{}
func (h MyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	switch strings.ToLower(r.URL.Path){
	case "/":
		fmt.Fprintln(w, "Welcome to home page")
	case "/login":
		handleLogin(w,r)
	case "/register":
		handleRegister(w,r)
	default:
		http.NotFound(w,r)
	}
}
var GlobalDbDataError bool = false
func readData(){
	file, err := os.ReadFile("data.json")
	if err != nil{
		GlobalDbDataError = true
		return
	}
	var db dbData
	err = json.Unmarshal(file, &db)
	if err != nil{
		GlobalDbDataError = true
	}
	
	
}
func main() {
	var h MyHandler
	fmt.Println("Listen and served")
	http.ListenAndServe(":8080", h)
	// psqlinfo := fmt.Sprintf("host=%s port=%d user=%s "+
	// 	"password=%s dbname=%s sslmode=disable",
	// 	host, port, user, password, dbname)
	// // open a connection
	// db, err := sql.Open("postgres", psqlinfo)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()
	// // check the connection
	// err = db.Ping()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// rows, err2 := db.Query("select * from users")
	// if err2 != nil {
	// 	log.Fatal(err2)
	// }
	// defer rows.Close()
	// var count int=0
	// for rows.Next() {
	// 	count++
	// 	var (
	// 		userid     int
	// 		createdon  string // assuming created_on is a date or timestamp
	// 		name       string
	// 	)
	// 	err2 := rows.Scan(&userid, &createdon, &name)
	// 	if err2 != nil {
	// 		log.Fatal(err2)
	// 	}
	// 	fmt.Printf("user_id = %d, created_on = %s, name = %s\n", userid, createdon, name)
	// }
	// // 	{stmt := `insert into users ( name) values ($1)`
	// // 	res, err2 := db.exec(stmt,"john doe") // replace with your values
	// // 	if err2 != nil {
	// // 		log.Fatal(err2)
	// // 	}
	// // fmt.Println(res)
	// // fmt.Println("row inserted successfully.")
	// // 	}
	// fmt.Printf("%d is the number of items \n", count)
	// // check for errors from iterating over rows.
	// err2 = rows.Err()
	// if err2 != nil {
	// 	log.Fatal(err2)
	// }
	// fmt.Println("successfully connected!")
}
