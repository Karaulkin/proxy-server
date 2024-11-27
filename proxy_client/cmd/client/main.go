package main

import (
	"context"
	"flag"
	"fmt"
	thumbnailv1 "github.com/Karaulkin/proto_thumbnail/gen/go/thumbnail"
	"google.golang.org/grpc"
	"log/slog"
	"os"
	"proxy_client/lib/logger/handlers/slogpretty"
	"strings"
	"sync"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	urls := flag.String("urls", "", "Список ссылок на YouTube через запятую")
	async := flag.Bool("async", false, "Асинхронная загрузка")
	flag.Parse()

	log := setupLogger(envLocal)

	log.Info("starting client")

	//TODO: Используем имя сервиса proxy_server вместо localhost grpc.Dial("proxy_server:44044", grpc.WithInsecure())
	// TODO:чтобы клиент ждал подъем сервиса
	conn, err := grpc.Dial("localhost:44044", grpc.WithInsecure())
	if err != nil {
		log.Error("failed to connect to server", slog.String("error", err.Error()))
		os.Exit(1) // Завершаем работу при невозможности подключиться
	}
	defer conn.Close()

	client := thumbnailv1.NewThumbnailClient(conn)
	links := strings.Split(*urls, ",")

	Run(log, client, links, async)
	log.Info("close client")
}

func Run(log *slog.Logger, client thumbnailv1.ThumbnailClient, links []string, async *bool) {
	if *async {
		var wg sync.WaitGroup
		for _, url := range links {
			wg.Add(1)
			go func(u string) {
				defer wg.Done()
				fetchThumbnail(client, log, u)
			}(url)
		}
		wg.Wait()
	} else {
		for _, url := range links {
			fetchThumbnail(client, log, url)
		}
	}
}

func fetchThumbnail(client thumbnailv1.ThumbnailClient, log *slog.Logger, url string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Hour)
	defer cancel()

	resp, err := client.GetThumbnail(ctx, &thumbnailv1.ThumbnailRequest{VideoUrl: url})
	if err != nil {
		log.Error("fetching thumbnail", slog.String("url", url), slog.String("error", err.Error()))
		return
	}

	log.Info("fetched thumbnail", slog.String("url", resp.VideoUrl), slog.Int("bytes", len(resp.ImageData)))
}

func waitForServer(log *slog.Logger, address string, timeout time.Duration, retryInterval time.Duration) (*grpc.ClientConn, error) {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(retryInterval))
		if err == nil {
			log.Info("Successfully connected to server", slog.String("address", address))
			return conn, nil
		}
		log.Warn("Server not available, retrying...", slog.String("address", address), slog.String("error", err.Error()))
		time.Sleep(retryInterval)
	}
	return nil, fmt.Errorf("server did not become available within %s", timeout)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
