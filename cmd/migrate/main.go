package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/sqlite3"
	"github.com/golang-migrate/migrate/source/file"
)

// File open ki instance bnaya db ka and then migrate krdiya
func main() {
	if len(os.Args) < 2 {
		log.Fatal("Provide a migration direction: 'Up' or 'Down'")
	}

	direction := os.Args[1]                     //Direction fetch ki migration ki
	db, err := sql.Open("sqlite3", "./data.db") //Database Connection open kiya
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close() //After completion database close kiya

	instance, err := sqlite3.WithInstance(db, &sqlite3.Config{}) //

	if err != nil {
		log.Fatal(err)
	}
	fileMig, err := (&file.File{}).Open("cmd/migrate/migrations") //Migrations wali file open ki migrations mein use krne ke liye
	// fileMig, err := os.Open("cmd/migrate/migrations")
	if err != nil {
		fmt.Println("Error Location Migrations file")
	}

	m, err := migrate.NewWithInstance("file", fileMig, "sqlite3", instance) //Migrate ka instance jismein File send krdi

	if err != nil {
		log.Fatal(err)
	}

	switch direction { //Direction check krke no change and other errors check krke Migration run krdo
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	default:
		log.Fatal("Invalid Direction")
	}
}
