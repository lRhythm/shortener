package rest

import "github.com/lRhythm/shortener/internal/models"

const (
	pathParamID = "id"
)

type createRequest struct {
	OriginalURL string `json:"url"`
}

type createItemRequest struct {
	OriginalURL   string `json:"original_url"`
	CorrelationID string `json:"correlation_id"`
}

type createItemsRequest []createItemRequest

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
