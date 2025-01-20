package grpc

import "github.com/lRhythm/shortener/internal/models"

// Функция convertURLCreateBatchRequestItemsToRows - преобразование списка из запроса метода пакетного создания URL
// в список для сервисного слоя.
func convertURLCreateBatchRequestItemsToRows(URLs []*LinkCreateBatchRequest_LinkItem) models.Rows {
	rows := make(models.Rows, 0, len(URLs))
	for _, item := range URLs {
		if len(item.OriginalLink) == 0 || len(item.CorrelationId) == 0 {
			continue
		}
		rows = append(rows, models.Row{
			OriginalURL:   item.OriginalLink,
			CorrelationID: item.CorrelationId,
		})
	}
	return rows
}

// Функция convertRowsToURLCreateBatchResponseItems - преобразование списка URL из сервисного слоя в список URL
// для ответа метода пакетного создания URL.
func convertRowsToURLCreateBatchResponseItems(rows models.Rows) []*LinkCreateBatchResponse_LinkItem {
	items := make([]*LinkCreateBatchResponse_LinkItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, &LinkCreateBatchResponse_LinkItem{
			ShortLink:     row.ShortURL,
			CorrelationId: row.CorrelationID,
		})
	}
	return items
}

// Функция convertRowsToUserURListResponseItems - преобразование списка URL из сервисного слоя в список URL
// для ответа метода пакетного получения URL пользователя.
func convertRowsToUserURListResponseItems(rows models.Rows) []*UserLinkListResponse_LinkItem {
	items := make([]*UserLinkListResponse_LinkItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, &UserLinkListResponse_LinkItem{
			ShortLink:    row.ShortURL,
			OriginalLink: row.OriginalURL,
		})
	}
	return items
}
