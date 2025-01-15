package grpc

import "github.com/lRhythm/shortener/internal/models"

// Функция convertURLCreateBatchRequestItemsToRows - преобразование списка из запроса метода пакетного создания URL
// в список для сервисного слоя.
func convertURLCreateBatchRequestItemsToRows(URLs []*UrlCreateBatchRequest_UrlItem) models.Rows {
	rows := make(models.Rows, 0, len(URLs))
	for _, item := range URLs {
		if len(item.OriginalUrl) == 0 || len(item.CorrelationId) == 0 {
			continue
		}
		rows = append(rows, models.Row{
			OriginalURL:   item.OriginalUrl,
			CorrelationID: item.CorrelationId,
		})
	}
	return rows
}

// Функция convertRowsToURLCreateBatchResponseItems - преобразование списка URL из сервисного слоя в список URL
// для ответа метода пакетного создания URL.
func convertRowsToURLCreateBatchResponseItems(rows models.Rows) []*UrlCreateBatchResponse_UrlItem {
	items := make([]*UrlCreateBatchResponse_UrlItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, &UrlCreateBatchResponse_UrlItem{
			ShortUrl:      row.ShortURL,
			CorrelationId: row.CorrelationID,
		})
	}
	return items
}

// Функция convertRowsToUserURListResponseItems - преобразование списка URL из сервисного слоя в список URL
// для ответа метода пакетного получения URL пользователя.
func convertRowsToUserURListResponseItems(rows models.Rows) []*UserUrlListResponse_UrlItem {
	items := make([]*UserUrlListResponse_UrlItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, &UserUrlListResponse_UrlItem{
			ShortUrl:    row.ShortURL,
			OriginalUrl: row.OriginalURL,
		})
	}
	return items
}
