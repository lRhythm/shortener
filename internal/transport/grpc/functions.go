package grpc

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// Метод trusted - проверка вхождения IP-адреса клиента в доверенную подсеть.
func (s *Server) trusted(ctx context.Context) error {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return status.Error(codes.Internal, "failed get peer info from ctx")
	}
	if t := s.cfg.Trusted(); len(t) == 0 || !strings.HasPrefix(p.Addr.String(), t) {
		// Forbidden.
		return status.Error(codes.PermissionDenied, "permission denied")
	}
	return nil
}

// Метод userID - получение из metadata и проверка идентификатора пользователя (аутентификация).
func (s *Server) userID(ctx context.Context) (string, error) {
	var ID string
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ID, status.Error(codes.Internal, "failed get metadata from ctx")
	}
	values := md.Get(userID)
	if len(values) == 0 {
		// Unauthorized.
		return ID, status.Error(codes.Unauthenticated, "missing user_id")
	}
	ID = values[0]
	if err := s.service.ValidateUserID(ID); err != nil {
		// Unauthorized.
		return ID, status.Errorf(codes.Unauthenticated, "invalid user_id: %v", err)
	}
	return ID, nil
}

// Метод newUserID - создание идентификатора пользователя (регистрация).
// Если пользователь зарегистрирован, то значение идентификатора пользователя устанавливается в metadata ответа.
func (s *Server) newUserID(ctx context.Context) (string, error) {
	var ID string
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ID, status.Error(codes.Internal, "failed get metadata from ctx")
	}
	if values := md.Get(userID); len(values) > 0 {
		ID = values[0]
		if err := s.service.ValidateUserID(ID); err == nil {
			return ID, nil
		}
	}
	ID = s.service.GenerateUserID()
	// Задачу установки metadata можно решить с помощью unaryInterceptor, но:
	// - функционал регистрация будет распространяться на все gRPC методы;
	// - сервисный слой не доступен в unaryInterceptor -> необходимы публичные функциям вместо имплементаций для Server.
	md = metadata.New(map[string]string{userID: ID})
	err := grpc.SendHeader(ctx, md)
	if err != nil {
		return ID, status.Error(codes.Internal, "failed set metadata to header")
	}
	return ID, nil
}
