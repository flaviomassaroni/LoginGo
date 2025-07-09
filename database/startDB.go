package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type User struct {
	Username string `db:"username"`
	Password string `db:"password"`
}

func StartDB() *sqlx.DB {
	//connect to a PostgreSQL database
	// Replace the connection details (user, dbname, password, host) with your own by compose.yaml
	db, err := sqlx.Connect("postgres", "user=logingo dbname=mydb sslmode=disable password=mypassword host=db")
	if err != nil {
		log.Fatalln(err)
	}

	// Test the connection to the database
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully Connected")
	}

	// Crea tabella se non esiste
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY,
		password TEXT NOT NULL
	);`
	_, err = db.Exec(schema)
	if err != nil {
		log.Fatalf("Errore creazione tabella users: %v", err)
	} else {
		log.Println("Successfully Created users")
	}

	return db

}
