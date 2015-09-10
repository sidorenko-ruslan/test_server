package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"io"
	"strings"
	"bytes"
)

var (
	db *sql.DB
	err error

)
const (
	InsertUserQuery string = "INSERT INTO users (first_name,last_name,patronymic) VALUES ($1,$2,$3)"
	UpdateUserQuery string = "UPDATE users SET first_name = $1, last_name = $2, patronymic = $3 WHERE id=$4"
	DeleteUserQuery string = "DELETE FROM users WHERE id=$1"
	SelectUsersQuery string = "SELECT first_name, last_name, patronymic FROM users"
)

func handleGetQuery(w *http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	var rows *sql.Rows
	if len(id) > 0 {
		rows,err = db.Query(SelectUsersQuery + " WHERE id=$1",id)
	} else {
		rows,err = db.Query(SelectUsersQuery)
	}
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	var firstName string
	var lastName string
	var patronymic string

	var responseBuffer bytes.Buffer
	for rows.Next() {
		err := rows.Scan(&firstName,&lastName,&patronymic)
		if err != nil {
			log.Fatal(err)
		}
		responseBuffer.WriteString("user {first name = '" + firstName + "', last name = '" + lastName + "', patronymic = '" + patronymic + "'}\n")
	}
	if responseBuffer.Len() == 0 {
		responseBuffer.WriteString("no users\n")
	} 
	io.WriteString(*w, responseBuffer.String())
}

func handlePostQuery(w *http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	firstName := r.Form.Get("first_name")
	lastName := r.Form.Get("last_name")
	patronymic := r.Form.Get("patronymic")
	_,err = db.Exec(InsertUserQuery,firstName,lastName,patronymic)
	if err != nil {
		log.Fatal(err)
	}
	io.WriteString(*w, "the user with first name '" + firstName + "' has been successfully added\n")
}

func handlePutQuery(w *http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	r.ParseForm()
	firstName := r.Form.Get("first_name")
	lastName := r.Form.Get("last_name")
	patronymic := r.Form.Get("patronymic")
	_,err = db.Exec(UpdateUserQuery,firstName,lastName,patronymic,id)
	if err != nil {
		log.Fatal(err)
	}
	io.WriteString(*w, "the user with id '" + id + "' has been successfully changed\n")
}

func handleDeleteQuery(w *http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	_,err = db.Exec(DeleteUserQuery,id)
	if err != nil {
		log.Fatal(err)
	}
	io.WriteString(*w, "the user with id '" + id + "' has been successfully removed\n")
}

func handler(w http.ResponseWriter, r *http.Request) {
	if path := strings.Replace(r.URL.Path,"/","",-1); path == "users" {
		switch r.Method {
		case "GET":
			handleGetQuery(&w,r)
		case "POST":
			handlePostQuery(&w,r)
		case "PUT":
			handlePutQuery(&w,r)
		case "DELETE":
			handleDeleteQuery(&w,r)
		}
	}
}

func main() {
	db, err = sql.Open("postgres", "user=postgres dbname=test_users sslmode=disable")
	defer db.Close()
		if err != nil {
			log.Fatal(err)
	}
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080",nil)
}