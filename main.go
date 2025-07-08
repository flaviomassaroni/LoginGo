package main

import (
	"fmt"
	"goLogin/database"
	"net/http"
)

type Login struct {
	HashedPassword string
	SessionToken   string
	CSRFToken      string
}

func main() {

	// Starting up PostgreSLQ DB
	database.StartDB()

	// Api Endpoints
	http.HandleFunc("/", serveLoginPage)
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
	http.HandleFunc("/logout", logout)
	http.ListenAndServe(":8080", nil)
}

func serveLoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "template/home.html")
}

func register(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "template/register.html")
		return
	}
	if r.Method != http.MethodPost {
		err := http.StatusMethodNotAllowed
		http.Error(w, "Invalid method", err)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	hashedPassword, err := hashPassword(password)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//	_, err = database.DB.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, string(hashedPassword))
	//	if err != nil {
	//		http.Error(w, "Registration Error: "+err.Error(), http.StatusInternalServerError)
	//		return
	//	}

	fmt.Fprintln(w, "User registered successfully!", username, hashedPassword)

}

func login(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "template/login.html")
		return
	}

	if r.Method != http.MethodPost {
		err := http.StatusMethodNotAllowed
		http.Error(w, "Invalid method", err)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	fmt.Fprintln(w, "User logged successfully!", username, password)
}

func logout(w http.ResponseWriter, r *http.Request) {

}
