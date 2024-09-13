package storage

import (
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lRhythm/shortener/internal/models"
	"github.com/lib/pq"
)

type DB struct {
	db *sqlx.DB
}

func (db *DB) Ping() error {
	return db.db.Ping()
}

func (db *DB) Put(shortURL, originalURL string) error {
	query := `insert into urls (short_url, original_url) values ($1, $2)`
	_, err := db.db.Exec(query, shortURL, originalURL)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pgerrcode.UniqueViolation {
			err = models.ErrConflict
		}
	}
	return err
}

func (db *DB) Batch(rows models.Rows) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `insert into urls (short_url, original_url, correlation_id) values ($1, $2, $3)`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, row := range rows {
		_, err = stmt.Exec(row.ShortURL, row.OriginalURL, row.CorrelationID)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (db *DB) GetOriginalURL(shortURL string) (string, error) {
	var originalURL string
	query := `select original_url from urls where short_url = $1`
	err := db.db.Get(&originalURL, query, shortURL)
	return originalURL, err
}

func (db *DB) GetShortURL(originalURL string) (string, error) {
	var shortURL string
	query := `select short_url from urls where original_url = $1`
	err := db.db.Get(&shortURL, query, originalURL)
	return shortURL, err
}

func (db *DB) Close() error {
	return db.db.Close()
}

func NewDB(DSN string) (*DB, error) {
	db, err := sqlx.Open("postgres", DSN)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	err = migrate(db)
	if err != nil {
		return nil, err
	}
	return &DB{
		db: db,
	}, nil
}
