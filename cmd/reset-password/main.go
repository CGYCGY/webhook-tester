package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	email := flag.String("email", "", "email of the user to reset")
	password := flag.String("password", "", "new password (min 8 characters)")
	defaultDB := os.Getenv("DATA_DIR")
	if defaultDB == "" {
		defaultDB = "/data/webhook-tester.db"
	} else {
		defaultDB = defaultDB + "/webhook-tester.db"
	}
	dbPath := flag.String("db", defaultDB, "path to the SQLite database file")
	flag.Parse()

	if *email == "" {
		log.Fatal("--email is required")
	}
	if *password == "" {
		log.Fatal("--password is required")
	}
	if len(*password) < 8 {
		log.Fatal("password must be at least 8 characters")
	}

	db, err := sql.Open("sqlite3", *dbPath)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	hashed, err := bcrypt.GenerateFromPassword([]byte(*password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	result, err := db.Exec(
		"UPDATE users SET password = ?, updated_at = CURRENT_TIMESTAMP WHERE email = ?",
		string(hashed), *email,
	)
	if err != nil {
		log.Fatalf("failed to update password: %v", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatalf("failed to get rows affected: %v", err)
	}
	if rows == 0 {
		fmt.Fprintf(os.Stderr, "no user found with that email\n")
		os.Exit(1)
	}

	fmt.Printf("password reset successfully for %s\n", *email)
}
