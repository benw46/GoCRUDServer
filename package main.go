package main

import (
	"fmt"
	"log"
	"net/http"
	_"github.com/lib/pq"
	"github.com/jmoiron/sqlx"
)

const (
	host 		= "localhost"
	port 		= 5432
	user 		= "postgres"
	password 	= "password"
	dbname 		= "dbname"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to homepage")
	fmt.Println("Endpoint hit: homePage")
}

func handleRequests() {
	// Mux not working, http is working
	//myRouter := mux.NewRouter().StrictSlash(true)
	//myRouter.HandleFunc("/", homePage)
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func dbConnection() (*sqlx.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
    "password=%s dbname=%s sslmode=disable",
    host, port, user, password, dbname)

	db, err := sqlx.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
		fmt.Println("Error during sqlx.Open: ", err)
	}
	return db, err
}

func dbSetup(db *sqlx.DB, err error) {
	var schema = `
	CREATE TABLE place (
		country text PRIMARY KEY 
	);`
	db.Exec(schema)
}

func createEntry(db *sqlx.DB) {
	insertFunc := `INSERT INTO place (country) VALUES ('Aus')`
	db.Exec(insertFunc)

	insertFunc = `INSERT INTO place (country) VALUES ('Sgp')`
	db.Exec(insertFunc)
}

func readEntry(db *sqlx.DB) {
	rows, err := db.Query("SELECT country FROM place")
	if err != nil {
		fmt.Println("Error: ", err)
	}
	for rows.Next() {
		var country string
		err = rows.Scan(&country)
		fmt.Println(country)
	}
}

func updateEntry(db *sqlx.DB) {
	stmt := `UPDATE place SET country = 'Gbr' WHERE country = 'Aus';`
	_, err := db.Exec(stmt)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	readEntry(db)
}

func deleteEntry(db *sqlx.DB) {
	stmt := `DELETE FROM place WHERE country = 'Gbr'`
	_, err := db.Exec(stmt)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	readEntry(db)
}

func main () {
	var db *sqlx.DB
	var err error
	db, err = dbConnection()

	dbSetup(db, err)
	createEntry(db)
	updateEntry(db)
	deleteEntry(db)
    handleRequests()
}
