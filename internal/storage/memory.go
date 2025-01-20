package storage

import (
	"errors"
	"slices"

	"github.com/lRhythm/shortener/internal/models"
)

// Memory - объект пакета для взаимодействия с памятью.
type Memory struct {
	storage *[]Row
	file    *file
}

// Ping - заглушка для соответствия интерфейсу service.RepositoryInterface.
func (m *Memory) Ping() error {
	return nil
}

// CountURL - операция получения количества сокращённых URL в сервисе.
func (m *Memory) CountURL() (uint, error) {
	return uint(len(*m.storage)), nil
}

// CountUser - операция получения количества пользователей в сервисе.
func (m *Memory) CountUser() (uint, error) {
	cm := make(map[string]bool)
	for _, v := range *m.storage {
		if _, ok := cm[v.UserID]; !ok {
			cm[v.UserID] = true
		}
	}
	return uint(len(cm)), nil
}

// Put - операция добавления в память сокращенного URL.
func (m *Memory) Put(shortURL, originalURL, userID string) error {
	for _, row := range *m.storage {
		if originalURL == row.OriginalURL {
			return models.ErrConflict
		}
	}
	*m.storage = append(*m.storage, newRow(shortURL, originalURL, "", userID))
	return nil
}

// Batch - операция пакетного добавления в память сокращенного URL.
func (m *Memory) Batch(rows models.Rows, userID string) error {
	for _, row := range rows {
		*m.storage = append(*m.storage, newRow(row.ShortURL, row.OriginalURL, row.CorrelationID, userID))
	}
	return nil
}

// GetOriginalURL - операция получения из памяти исходного URL по сокращенному.
func (m *Memory) GetOriginalURL(shortURL string) (string, bool, error) {
	for _, row := range *m.storage {
		if shortURL == row.ShortURL {
			return row.OriginalURL, row.IsDeleted, nil
		}
	}
	return "", false, errors.New("short url not found")
}

// GetShortURL - операция добавления в память сокращенного URL по исходному.
func (m *Memory) GetShortURL(originalURL string) (string, error) {
	for _, row := range *m.storage {
		if originalURL == row.OriginalURL {
			return row.ShortURL, nil
		}
	}
	return "", errors.New("original url not found")
}

// GetUserURLs - операция получения из памяти сокращенных URL пользователя.
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

// DeleteUserURLS - операция удаления из памяти сокращенных URL пользователя.
func (m *Memory) DeleteUserURLS(shortURLs []string, userID string) error {
	for i, row := range *m.storage {
		if userID == row.UserID && slices.Contains(shortURLs, row.ShortURL) {
			row.IsDeleted = true
			(*m.storage)[i] = row
		}
	}
	return nil
}

// Close - закрытие клиента работы с памятью: запись в файл данных из хранилища в памяти.
func (m *Memory) Close() error {
	defer m.file.close()
	err := m.file.writeRows(m.storage)
	if err != nil {
		return err
	}
	return nil
}

// NewMemory - создание объекта хранилища в памяти.
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
