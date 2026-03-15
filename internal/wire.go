package wire

import (
	"github.com/google/wire"

	"user-center/internal/biz"
	"user-center/internal/conf"
	"user-center/internal/data"
	"user-center/internal/server"
	"user-center/internal/service"
)

func NewApp(conf *conf.Config, logger log.Logger) (*App, func(), error) {
	panic(wire.Build(
		conf.NewConfig,
		data.NewDB,
		data.NewCache,
		data.NewData,
		biz.NewAuthUseCase,
		service.NewUserService,
		server.NewHTTPServer,
		server.NewGRPCServer,
		NewApp,
	))
}

type App struct {
	http *http.Server
	grpc *grpc.Server
}

func (a *App) Run() error {
	return nil
}
