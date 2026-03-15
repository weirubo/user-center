package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"strings"

	v1 "user-center/api/user/v1"
	"user-center/internal/biz"
	"user-center/internal/conf"
	"user-center/internal/data"
	"user-center/internal/handler"
	"user-center/internal/server"
	"user-center/internal/service"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newApp(logger log.Logger, gs *grpc.Server) *kratos.App {
	return kratos.New(
		kratos.Name("user-center"),
		kratos.Logger(logger),
		kratos.Server(gs),
	)
}

func authMiddleware(authUC *biz.AuthUseCase, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "unauthorized", 401)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "unauthorized", 401)
			return
		}

		userID, err := authUC.ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "unauthorized", 401)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next(w, r.WithContext(ctx))
	}
}

func main() {
	flag.Parse()
	logger := log.NewStdLogger(os.Stdout)

	c := config.New(
		config.WithSource(file.NewSource("configs/config.yaml")),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Config
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	db, err := data.NewDB(bc.Database)
	if err != nil {
		panic(err)
	}
	cache, _ := data.NewCache(nil)
	defer func() {
		db.Close()
		if cache != nil {
			cache.Close()
		}
	}()

	userRepo := data.NewUserRepo(db, cache)
	authUC := biz.NewAuthUseCase(userRepo, bc.Auth.JwtSecret, bc.Auth.ExpireTime, bc.Smtp)
	userService := service.NewUserService(authUC)
	userHandler := handler.NewUserHandler(authUC)

	gs := server.NewGRPCServer(bc.Server, userService, logger)

	// Start HTTP gateway in goroutine
	go func() {
		ctx := context.Background()
		mux := runtime.NewServeMux()
		opts := []ggrpc.DialOption{
			ggrpc.WithTransportCredentials(insecure.NewCredentials()),
		}
		_ = v1.RegisterUserHandlerFromEndpoint(ctx, mux, "localhost:9000", opts)

		// Add custom HTTP routes using http.ServeMux
		httpMux := http.NewServeMux()
		httpMux.HandleFunc("/api/v1/register", userHandler.Register)
		httpMux.HandleFunc("/api/v1/login", userHandler.Login)
		httpMux.HandleFunc("/api/v1/userinfo", authMiddleware(authUC, userHandler.GetUserInfo))
		httpMux.HandleFunc("/api/v1/account/delete", authMiddleware(authUC, userHandler.DeleteAccount))
		httpMux.HandleFunc("/api/v1/account", userHandler.DeleteAccount)
		httpMux.HandleFunc("/api/v1/password/change", authMiddleware(authUC, userHandler.ChangePassword))
		httpMux.HandleFunc("/api/v1/verifycode/send", userHandler.SendVerifyCode)

		http.ListenAndServe(":8000", httpMux)
	}()

	app := newApp(logger, gs)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
