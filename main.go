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

	var secretKey []byte = []byte("secret")

	// Api Endpoints
	http.HandleFunc("/", ServeLoginPage)
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		Login(w, r, db, secretKey)
	})
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		Register(w, r, db)
	})
	http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
		Home(w, r, secretKey)
	})

	http.HandleFunc("/logout", Logout)
	http.ListenAndServe(":8080", nil)

}

func ServeLoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "template/temporary.html")
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

	if UsernameExists(username, db) {
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

func Login(w http.ResponseWriter, r *http.Request, db *sqlx.DB, secretKey []byte) {

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

	if !UsernameExists(username, db) {
		http.Error(w, "Username does not exists!", http.StatusUnauthorized)
		return
	}

	var toCheck string

	err := db.Get(&toCheck, "SELECT password FROM users WHERE username=($1)", username)
	if err != nil {
		http.Error(w, "Invalid Password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(toCheck), []byte(password))
	if err != nil {
		http.Error(w, "Invalid Password", http.StatusUnauthorized)
		return
	}
	tokenString, err := CreatingJWTToken(username, secretKey)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(map[string]string{"token": tokenString})

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		HttpOnly: true,
		Path:     "/",
	})

	http.Redirect(w, r, "/home", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {

}

func Home(w http.ResponseWriter, r *http.Request, secretKey []byte) {

	cookie, error := r.Cookie("token")
	if error != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err := VerifyToken(cookie.Value, secretKey)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	http.ServeFile(w, r, "template/home.html")
}
