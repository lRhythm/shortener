package storage

import "errors"

type InMemory struct {
	s map[string]string
}

func (m *InMemory) Put(key, value string) error {
	_, ok := m.s[key]
	if ok {
		return errors.New("key already exists")
	}
	m.s[key] = value
	return nil
}

func (m *InMemory) Get(key string) (string, error) {
	v, ok := m.s[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return v, nil
}

func NewInMemory() *InMemory {
	return &InMemory{
		s: map[string]string{},
	}
}
