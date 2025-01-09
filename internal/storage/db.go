package storage

import (
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/lRhythm/shortener/internal/models"
)

// DB - объект пакета для взаимодействия с БД.
type DB struct {
	db *sqlx.DB
}

// Ping - пинг БД для реализации healthcheck.
func (db *DB) Ping() error {
	return db.db.Ping()
}

// CountURL - выполнение запроса получения количества сокращённых URL в сервисе.
func (db *DB) CountURL() (uint, error) {
	var cnt uint
	query := `select count(id) from urls`
	err := db.db.Get(&cnt, query)
	return cnt, err
}

// CountUser - выполнение запроса получения количества пользователей в сервисе.
func (db *DB) CountUser() (uint, error) {
	var cnt uint
	query := `select count(distinct(user_id)) from urls`
	err := db.db.Get(&cnt, query)
	return cnt, err
}

// Put - выполнение запроса добавления в БД сокращенного URL.
func (db *DB) Put(shortURL, originalURL, userID string) error {
	query := `insert into urls (short_url, original_url, user_id) values ($1, $2, $3)`
	_, err := db.db.Exec(query, shortURL, originalURL, userID)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pgerrcode.UniqueViolation {
			err = models.ErrConflict
		}
	}
	return err
}

// Batch - выполнение запросов пакетного добавления в БД сокращенного URL.
func (db *DB) Batch(rows models.Rows, userID string) error {
	tx, err := db.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := `insert into urls (short_url, original_url, correlation_id, user_id) values ($1, $2, $3, $4)`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, row := range rows {
		_, err = stmt.Exec(row.ShortURL, row.OriginalURL, row.CorrelationID, userID)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

// GetOriginalURL - выполнение запроса получения из БД исходного URL по сокращенному.
func (db *DB) GetOriginalURL(shortURL string) (string, bool, error) {
	var row models.Row
	query := `select original_url, is_deleted from urls where short_url = $1`
	err := db.db.Get(&row, query, shortURL)
	return row.OriginalURL, row.IsDeleted, err
}

// GetShortURL - выполнение запроса добавление в БД сокращенного URL по исходному.
func (db *DB) GetShortURL(originalURL string) (string, error) {
	var shortURL string
	query := `select short_url from urls where original_url = $1`
	err := db.db.Get(&shortURL, query, originalURL)
	return shortURL, err
}

// GetUserURLs - выполнение запроса получения из БД сокращенных URL пользователя.
func (db *DB) GetUserURLs(userID string) (models.Rows, error) {
	rows := make(models.Rows, 0)
	query := `select short_url, original_url from urls where user_id = $1`
	err := db.db.Select(&rows, query, userID)
	return rows, err
}

// DeleteUserURLS - выполнение запроса удаления из БД сокращенных URL пользователя.
func (db *DB) DeleteUserURLS(shortURLs []string, userID string) error {
	query := `update urls set is_deleted = true where short_url = any($1) and user_id = $2`
	_, err := db.db.Exec(query, pq.Array(shortURLs), userID)
	return err
}

// Close - закрытие соединения БД.
func (db *DB) Close() error {
	return db.db.Close()
}

// NewDB - создание объекта БД.
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
