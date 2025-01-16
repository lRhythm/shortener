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
		Links: uint32(countURL),
	}, nil
}

// LinkCreate - метод создания сокращенного URL.
func (s *Server) LinkCreate(ctx context.Context, in *LinkCreateRequest) (*LinkCreateResponse, error) {
	user, err := s.newUserID(ctx)
	if err != nil {
		return nil, err
	}
	shortURL, err := s.service.CreateShortURL(in.Link, "", user)
	if err != nil {
		if errors.Is(err, models.ErrConflict) {
			// Conflict.
			return &LinkCreateResponse{
				Result: shortURL,
			}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &LinkCreateResponse{
		Result: shortURL,
	}, nil
}

// LinkCreateBatch - метод пакетного создания сокращенных URL.
func (s *Server) LinkCreateBatch(ctx context.Context, in *LinkCreateBatchRequest) (*LinkCreateBatchResponse, error) {
	user, err := s.newUserID(ctx)
	if err != nil {
		return nil, err
	}
	rows, err := s.service.CreateBatch(convertURLCreateBatchRequestItemsToRows(in.Links), "", user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &LinkCreateBatchResponse{
		Links: convertRowsToURLCreateBatchResponseItems(rows),
	}, nil
}

// LinkGet - метод получения оригинального URL по сокращенному URL.
func (s *Server) LinkGet(ctx context.Context, in *LinkGetRequest) (*LinkGetResponse, error) {
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
		// Подобная реализация представлена в LinkCreate.
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
	return &LinkGetResponse{
		Link: originalURL,
	}, nil
}

// UserLinkList - метод получения сокращенных URL пользователя.
func (s *Server) UserLinkList(ctx context.Context, in *emptypb.Empty) (*UserLinkListResponse, error) {
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
	return &UserLinkListResponse{
		Links: convertRowsToUserURListResponseItems(rows),
	}, nil
}

// UserLinkDelete - метод удаления сокращенных URL пользователя.
func (s *Server) UserLinkDelete(ctx context.Context, in *UserLinkDeleteRequest) (*emptypb.Empty, error) {
	user, err := s.userID(ctx)
	if err != nil {
		return nil, err
	}
	go s.service.DeleteUserURLs(in.Links, user)
	// Accepted.
	return nil, nil
}
