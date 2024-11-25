package rest

import "github.com/lRhythm/shortener/internal/models"

const (
	pathParamID = "id"
)

// createRequest - request создания сокращенного URL.
type createRequest struct {
	OriginalURL string `json:"url"`
}

// createItemRequest - элемент коллекции request'а создания сокращенного URL.
type createItemRequest struct {
	OriginalURL   string `json:"original_url"`
	CorrelationID string `json:"correlation_id"`
}

// createItemsRequest - request создания сокращенного URL в виде коллекции элементов.
type createItemsRequest []createItemRequest

// ToRows - преобразование *createItemsRequest в models.Rows.
func (items *createItemsRequest) ToRows() models.Rows {
	var rows models.Rows
	for _, item := range *items {
		if item.OriginalURL == "" || item.CorrelationID == "" {
			continue
		}
		rows = append(rows, models.Row{
			OriginalURL:   item.OriginalURL,
			CorrelationID: item.CorrelationID,
		})
	}
	return rows
}
