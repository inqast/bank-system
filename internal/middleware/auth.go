package middleware

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type (
	parser interface {
		ParseJWT(authHeader []string) string
	}
)

func NewJWTInterceptor(parser parser) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return handler(ctx, req)
		}

		authHeader, ok := md["authorization"]
		if !ok || len(authHeader) == 0 {
			return handler(ctx, req)
		}

		userID := parser.ParseJWT(authHeader)

		newCtx := context.WithValue(ctx, "userID", userID)
		return handler(newCtx, req)
	}
}
