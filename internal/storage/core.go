package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq"
)

var schema = `
CREATE TABLE IF NOT EXISTS search (
    input text,
    output text
);
`

type storageImpl struct {
	db *sqlx.DB
}

type Storage interface {
	SetUp()
	GetOutput(input string) (string, error)
}

func NewStorage() Storage {
	return &storageImpl{}
}

func (s *storageImpl) SetUp() {
	db, err := sqlx.Connect("postgres", "user=postgres dbname=postgres password=postgres sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	db.MustExec(schema)

	tx := db.MustBegin()
	tx.MustExec("INSERT INTO search (input, output) VALUES ($1, $2)", "ТОРЖЕСТВОВАТЬ ПОБЕДА", "ГРОМКО") // КАК ТОРЖЕСТВОВАТЬ ПОБЕДУ?
	tx.Commit()
	s.db = db
}

func (s *storageImpl) GetOutput(input string) (string, error) {
	var output string

	err := s.db.Get(&output, "SELECT output FROM search WHERE TRIM(input) ILIKE $1 LIMIT 1;", input)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "Ничего не найдо", nil
		}
		return "", fmt.Errorf("ошибка при получении вывода: %v", err)
	}
	return output, nil
}
