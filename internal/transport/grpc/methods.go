package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/lRhythm/shortener/internal/models"
)

// Ping - метод ping.
func (s *Server) Ping(ctx context.Context, in *emptypb.Empty) (*emptypb.Empty, error) {
	if err := s.service.Ping(); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return nil, nil
}

// InternalStats - метод получения статистики: кол-ва сокращенных URL и кол-ва пользователей.
func (s *Server) InternalStats(ctx context.Context, in *emptypb.Empty) (*InternalStatsResponse, error) {
	err := s.trusted(ctx)
	if err != nil {
		return nil, err
	}
	countURL, err := s.service.CountURL()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	countUser, err := s.service.CountUser()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &InternalStatsResponse{
		Users: uint32(countUser),
		Urls:  uint32(countURL),
	}, nil
}

// UrlCreate - метод создания сокращенного URL.
func (s *Server) UrlCreate(ctx context.Context, in *UrlCreateRequest) (*UrlCreateResponse, error) {
	user, err := s.newUserID(ctx)
	if err != nil {
		return nil, err
	}
	shortURL, err := s.service.CreateShortURL(in.Url, "", user)
	if err != nil {
		if errors.Is(err, models.ErrConflict) {
			// Conflict.
			return &UrlCreateResponse{
				Result: shortURL,
			}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &UrlCreateResponse{
		Result: shortURL,
	}, nil
}

// UrlCreateBatch - метод пакетного создания сокращенных URL.
func (s *Server) UrlCreateBatch(ctx context.Context, in *UrlCreateBatchRequest) (*UrlCreateBatchResponse, error) {
	user, err := s.newUserID(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := s.service.CreateBatch(convertURLCreateBatchRequestItemsToRows(in.Urls), "", user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &UrlCreateBatchResponse{
		Urls: convertRowsToURLCreateBatchResponseItems(rows),
	}, nil
}

// UrlGet - метод получения оригинального URL по сокращенному URL.
func (s *Server) UrlGet(ctx context.Context, in *UrlGetRequest) (*UrlGetResponse, error) {
	if len(in.Id) == 0 {
		// Bad Request.
		return nil, status.Error(codes.InvalidArgument, "missing or empty id")
	}
	originalURL, isDeleted, err := s.service.GetOriginalURL(in.Id)
	if err != nil {
		// В пакете rest данная ошибка не обрабатывается намеренно, т.к. в постановке инкремента указано,
		// что при возникновении любой ошибки должен возвращаться ответ Bad Request.
		// Одним из решений будет создание пользовательской ошибки (например models.NoRows)
		// и возврат из storage слоя при её возникновении, например вместо sql.ErrNoRows.
		// Далее данную ошибку обработать как Not Found.
		// Подобная реализация представлена в UrlCreate.
		//if errors.Is(err, models.ErrNoRows) {
		//	// Not Found.
		//	return nil, status.Error(codes.NotFound, "urls not found")
		//}
		return nil, status.Error(codes.Internal, err.Error())
	}
	if isDeleted {
		// Gone.
		return nil, status.Error(codes.NotFound, "url is deleted")
	}
	return &UrlGetResponse{
		Url: originalURL,
	}, nil
}

// UserUrlList - метод получения сокращенных URL пользователя.
func (s *Server) UserUrlList(ctx context.Context, in *emptypb.Empty) (*UserUrlListResponse, error) {
	user, err := s.userID(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := s.service.GetUserURLs("", user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if len(rows) == 0 {
		// NoContent.
		return nil, nil
	}
	return &UserUrlListResponse{
		Urls: convertRowsToUserURListResponseItems(rows),
	}, nil
}

// UserUrlDelete - метод удаления сокращенных URL пользователя.
func (s *Server) UserUrlDelete(ctx context.Context, in *UserUrlDeleteRequest) (*emptypb.Empty, error) {
	user, err := s.userID(ctx)
	if err != nil {
		return nil, err
	}
	go s.service.DeleteUserURLs(in.Urls, user)
	// Accepted.
	return nil, nil
}
