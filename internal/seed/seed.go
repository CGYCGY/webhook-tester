package seed

import (
	"database/sql"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin(db *sql.DB, email, password string) error {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		return err
	}

	if count > 0 {
		log.Println("admin user already exists, skipping seed")
		return nil
	}

	id := uuid.New().String()
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"INSERT INTO users (id, email, password, created_at, updated_at) VALUES (?, ?, ?, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)",
		id, email, string(hashed),
	)
	if err != nil {
		return err
	}

	log.Println("admin user seeded successfully")
	return nil
}
