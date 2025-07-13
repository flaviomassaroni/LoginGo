package main

import (
	"fmt"
	"goLogin/database"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/jmoiron/sqlx"
)

//type login struct {
//	HashedPassword string
//	SessionToken   string
//	CSRFToken      string
//}

func main() {

	// Starting up PostgreSLQ DB
	db := database.StartDB()
	defer db.Close()

	// Api Endpoints
	http.HandleFunc("/", ServeLoginPage)
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		Login(w, r, db)
	})
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		Register(w, r, db)
	})
	http.HandleFunc("/logout", Logout)
	http.ListenAndServe(":8080", nil)
}

func ServeLoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "template/home.html")
}

func Register(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {

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

	if !CheckUsername(username) {
		http.Error(w, "Invalid Username", http.StatusInternalServerError)
		return
	}

	if UsernameExists(w, username, db) {
		http.Error(w, "Username already exists", http.StatusInternalServerError)
		return
	}

	if !CheckPassword(password) {
		http.Error(w, "Invalid Password", http.StatusInternalServerError)
		return
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", username, string(hashedPassword))
	if err != nil {
		http.Error(w, "Registration Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "User registered successfully!")

}

func Login(w http.ResponseWriter, r *http.Request, db *sqlx.DB) {

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

	if UsernameExists(w, username, db) {
		toCheck := fmt.Sprintf("SELECT password FROM users WHERE username=($1)", username)
		err := bcrypt.CompareHashAndPassword([]byte(toCheck), []byte(password))
		if err != nil {
			http.Error(w, "Invalid Password", http.StatusUnauthorized)
			return
		} else {
			fmt.Fprintln(w, "Login Successful!")
			// Set JWT token and redirect to home page

			// Generate JWT token
			token := CreatingJWTToken(username)

			return
		}

	}

}

func Logout(w http.ResponseWriter, r *http.Request) {

}
