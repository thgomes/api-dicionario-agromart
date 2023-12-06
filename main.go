package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

type Csa struct {
	ID             int    `json:"id"`
	NomeCSA        string `json:"nomeCSA"`
	ResponsavelCSA string `json:"responsavelCSA"`
	EmailCSA       string `json:"emailCSA"`
	URLBase        string `json:"urlBase"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./csas.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	createTable()

	http.HandleFunc("/csas", ListAllCSAs)
	http.HandleFunc("/csas/", GetCSA)
	http.HandleFunc("/csas/create", CreateCSA)

	fmt.Println("Server started on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func createTable() {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS csas (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nomeCSA TEXT NOT NULL,
		responsavelCSA TEXT NOT NULL,
		emailCSA TEXT NOT NULL,
		urlBase TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}
}

func ListAllCSAs(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM csas")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	csas := []Csa{}
	for rows.Next() {
		csa := Csa{}
		err := rows.Scan(&csa.ID, &csa.NomeCSA, &csa.ResponsavelCSA, &csa.EmailCSA, &csa.URLBase)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		csas = append(csas, csa)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(csas)
}

func GetCSA(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/csas/"):]
	row := db.QueryRow("SELECT * FROM csas WHERE id = ?", id)

	csa := Csa{}
	err := row.Scan(&csa.ID, &csa.NomeCSA, &csa.ResponsavelCSA, &csa.EmailCSA, &csa.URLBase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(csa)
}

func CreateCSA(w http.ResponseWriter, r *http.Request) {
	var csa Csa
	err := json.NewDecoder(r.Body).Decode(&csa)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO csas (nomeCSA, responsavelCSA, emailCSA, urlBase) VALUES (?, ?, ?, ?)",
		csa.NomeCSA, csa.ResponsavelCSA, csa.EmailCSA, csa.URLBase)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastID, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	csa.ID = int(lastID)

	w.WriteHeader(http.StatusCreated)
}
