package middleware

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"user-center/internal/biz"
)

func JWTAuth(jwtSecret string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, in interface{}) (interface{}, error) {
			if tr, ok := transport.FromServerContext(ctx); ok {
				// Get full path
				path := tr.RequestHeader().Get(":path")

				// Skip auth for register and login
				if path == "/api.v1.User/Register" || path == "/api.v1.User/Login" {
					return handler(ctx, in)
				}

				authHeader := tr.RequestHeader().Get("Authorization")
				if authHeader == "" {
					return nil, status.Error(codes.Unauthenticated, "missing authorization header")
				}

				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) != 2 || parts[0] != "Bearer" {
					return nil, status.Error(codes.Unauthenticated, "invalid authorization header format")
				}
				token := parts[1]

				uc := biz.NewAuthUseCase(nil, jwtSecret, 0, nil)
				userID, err := uc.ValidateToken(token)
				if err != nil {
					return nil, status.Error(codes.Unauthenticated, "invalid token")
				}

				ctx = context.WithValue(ctx, "user_id", userID)
			}
			return handler(ctx, in)
		}
	}
}
