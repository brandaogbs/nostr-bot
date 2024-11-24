package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
}

func InitDB(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	schemaPath := filepath.Join("database", "schema.sql")
	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(string(schema))
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) IsContentStored(contentID string) bool {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM content WHERE content_id = ?)`
	err := db.QueryRow(query, contentID).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if content is stored: %v", err)
		return false
	}
	return exists
}

func (db *DB) InsertRetrievedContent(contentID, content, source string) error {
	_, err := db.Exec(`INSERT INTO content(content_id, content, source) VALUES (?, ?, ?)`,
		contentID, content, source)
	return err
}

// func CloseDB() {
// 	if db != nil {
// 		err := db.Close()
// 		if err != nil {
// 			log.Printf("Error closing the database: %v", err)
// 		}
// 	}
// }
