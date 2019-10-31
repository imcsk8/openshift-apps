package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

// Global variables

var username, password, database string
var DB *sql.DB

func main() {
	username, _ = os.LookupEnv("MYSQL_USER")
	password, _ = os.LookupEnv("MYSQL_PASSWORD")
	database, _ = os.LookupEnv("MYSQL_DATABASE")
	DB = db_connect(username, password, database)
	defer DB.Close()

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/setup", SetupHandler)
	r.HandleFunc("/testdata", TestDataHandler)
	r.HandleFunc("/testget", TestGetHandler)
	http.Handle("/", r)
	log.Printf("Starting Application\nServices:\n/\n/setup\n/testdata\n/testget")
	log.Fatal(http.ListenAndServe(":8000", r))
}

// Handlers

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	err := DB.Ping()
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Database Setup Variables:\nMYSQL_USER: %s\nMYSQL_PASSWORD: %s\nMYSQL_DATABASE: %s\n", username, password, database)
}

func SetupHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Creating schema")
	db_create_schema()
	fmt.Fprintf(w, "Database Setup Completed\n")
}

func TestDataHandler(w http.ResponseWriter, r *http.Request) {
	stmtIns, err := DB.Prepare("INSERT INTO test_mysql VALUES( ?, ? )")
	if err != nil {
		panic(err.Error())
	}
	stmtIns.Exec(1000, "test_name")
	fmt.Fprintf(w, "Test data added\n")
	stmtIns.Close()

}

func TestGetHandler(w http.ResponseWriter, r *http.Request) {
	var id int
	var name string
	err := DB.QueryRow("SELECT id, name FROM test_mysql").Scan(&id, &name)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Test Database Values:\nid: %d\nusername: %s\n", id, name)
}

// Data functions

func db_create_schema() {
	crt, err := DB.Prepare("CREATE TABLE test_mysql (id int, name varchar(20))")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	crt.Exec()
	crt.Close() // Close the statement when we leave main() / the program terminates

}

func db_connect(username string, password string, database string) *sql.DB {
	connstring := username + ":" + password + "@/" + database
	log.Print("Connecting to the database: " + connstring)
	db, err := sql.Open("mysql", connstring)
	if err != nil {
		panic(err.Error())
	}
	return db
}
