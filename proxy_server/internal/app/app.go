package app

import (
	"log/slog"
	grpcapp "proxy_server/internal/app/grpc"
	"proxy_server/internal/services/proxy"
	"proxy_server/internal/storage/cache"
)

// Обёртка для моего grpc app в виде запуска стопа и инициализации

type App struct {
	GRPCServer   *grpcapp.App
	CacheStorage *cache.Storage
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
) *App {
	storage, err := cache.New(storagePath)
	if err != nil {
		panic(err)
	}

	proxyService := proxy.New(log, storage, storage)

	grpcApp := grpcapp.New(log, proxyService, grpcPort)

	return &App{
		GRPCServer:   grpcApp,
		CacheStorage: storage,
	}
}
