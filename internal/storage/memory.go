package storage

type InMemory struct {
	s map[string]string
}

func (m *InMemory) Put(key, value string) {
	m.s[key] = value
}

func (m *InMemory) Get(key string) string {
	return m.s[key]
}

func NewInMemory() *InMemory {
	return &InMemory{
		s: map[string]string{},
	}
}
