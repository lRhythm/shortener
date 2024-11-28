package rest

import (
	"encoding/json"
	"testing"

	gojson "github.com/goccy/go-json"

	"github.com/lRhythm/shortener/internal/models"
)

// mock - генерация моковых данных для BenchmarkEncode и BenchmarkDecode.
func mock() []models.Row {
	size := 1_000
	s := make([]models.Row, size)
	for i := 0; i < size; i++ {
		s[i] = models.Row{
			ShortURL:      "ShortURL",
			OriginalURL:   "http://example.com",
			CorrelationID: "00000000-0000-0000-0000-000000000000",
			IsDeleted:     false,
		}
	}
	return s
}

// BenchmarkEncode - проверка encode пакетов encoding/json и goccy/go-json.
func BenchmarkEncode(b *testing.B) {
	data := mock()
	b.ResetTimer()
	b.Run("(encoding/json).Marshal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = json.Marshal(data)
		}
	})
	b.Run("(goccy/go-json).Marshal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = gojson.Marshal(data)
		}
	})
}

// BenchmarkDecode - проверка decode пакетов encoding/json и goccy/go-json.
func BenchmarkDecode(b *testing.B) {
	data, _ := json.Marshal(mock())
	b.ResetTimer()
	b.Run("(encoding/json).Unmarshal", func(b *testing.B) {
		var tmp []models.Row
		for i := 0; i < b.N; i++ {
			_ = json.Unmarshal(data, &tmp)
		}
	})
	b.Run("(goccy/go-json).Unmarshal", func(b *testing.B) {
		var tmp []models.Row
		for i := 0; i < b.N; i++ {
			_ = gojson.Unmarshal(data, &tmp)
		}
	})
}
