package database

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type User struct {
	username string `db:"username"`
	password string `db:"password"`
}

func StartDB() {
	//connect to a PostgreSQL database
	// Replace the connection details (user, dbname, password, host) with your own by compose.yaml
	db, err := sqlx.Connect("postgres", "user=logingo dbname=yourdb sslmode=disable password=mypassword host=db")
	if err != nil {
		log.Fatalln(err)
	}

	defer db.Close()

	// Test the connection to the database
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	} else {
		log.Println("Successfully Connected")
	}

}
