package storage

import (
	"errors"
	"github.com/lRhythm/shortener/internal/models"
)

type Memory struct {
	storage *[]Row
	file    *file
}

func (m *Memory) Ping() error {
	return nil
}

func (m *Memory) Put(shortURL, originalURL, userID string) error {
	*m.storage = append(*m.storage, newRow(shortURL, originalURL, "", userID))
	return nil
}

func (m *Memory) Batch(rows models.Rows, userID string) error {
	for _, row := range rows {
		*m.storage = append(*m.storage, newRow(row.ShortURL, row.OriginalURL, row.CorrelationID, userID))
	}
	return nil
}

func (m *Memory) GetOriginalURL(shortURL string) (string, error) {
	for _, row := range *m.storage {
		if shortURL == row.ShortURL {
			return row.OriginalURL, nil
		}
	}
	return "", errors.New("short url not found")
}

func (m *Memory) GetShortURL(originalURL string) (string, error) {
	for _, row := range *m.storage {
		if originalURL == row.OriginalURL {
			return row.ShortURL, nil
		}
	}
	return "", errors.New("original url not found")
}

func (m *Memory) GetUserURLs(userID string) (models.Rows, error) {
	rows := make(models.Rows, 0)
	for _, row := range *m.storage {
		if userID == row.UserID {
			rows = append(rows, models.Row{
				ShortURL:    row.ShortURL,
				OriginalURL: row.OriginalURL,
			})
		}
	}
	return rows, nil
}

func (m *Memory) Close() error {
	defer m.file.close()
	err := m.file.writeRows(m.storage)
	if err != nil {
		return err
	}
	return nil
}

func NewMemory(fname string) (*Memory, error) {
	f, err := newFile(fname)
	if err != nil {
		return nil, err
	}
	s, err := f.readRows()
	if err != nil {
		return nil, err
	}
	return &Memory{
		storage: s,
		file:    f,
	}, nil
}
