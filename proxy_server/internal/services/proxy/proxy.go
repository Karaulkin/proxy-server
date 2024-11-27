package proxy

import (
	"context"
	"log/slog"
	"proxy_server/internal/domain/models"
	"proxy_server/internal/lib/downloder"
)

type Thumbnail struct {
	log           *slog.Logger
	cacheProvider CacheProvider
	cacheSaver    CacheSaver
	preview       *downloder.Downloader
}

type CacheProvider interface {
	CacheUrlProvider(ctx context.Context, url string) (models.CachedData, error)
}

type CacheSaver interface {
	SaveCache(ctx context.Context, url string, data []byte) (int, error)
}

func New(
	log *slog.Logger,
	cacheProvider CacheProvider,
	cacheProvide CacheSaver,
) *Thumbnail {
	return &Thumbnail{
		log:           log,
		cacheProvider: cacheProvider,
		cacheSaver:    cacheProvide,
		preview:       downloder.New(log),
	}
}

// Url check if url exists in the system and return cache data.
// If url doesn't exist, return error.
func (p *Thumbnail) Url(ctx context.Context, url string) ([]byte, string, error) {
	const op = "Thumbnail.UrlProvider"

	log := p.log.With(
		slog.String("op", op),
		slog.String("url", url),
	)

	data, err := p.cacheProvider.CacheUrlProvider(ctx, url)
	
	if err != nil {
		log.Info("downloading to url")

		img, err := p.preview.Download(url)
		if err != nil {
			log.Error("failed to download image preview")
			return nil, "", err
		}

		id, err := p.cacheSaver.SaveCache(ctx, url, img)
		if err != nil {
			log.Error("failed to save image to cache")
			return nil, "", err
		}

		log.Info("Cache image doesn't saved", slog.Any("id", id))
		return img, url, nil

	} else if data.Data == nil {
		log.Info("downloading to url")

		img, err := p.preview.Download(url)
		if err != nil {
			log.Error("failed to download image preview")
			return nil, "", err
		}

		id, err := p.cacheSaver.SaveCache(ctx, url, img)
		if err != nil {
			log.Error("failed to save image to cache")
			return nil, "", err
		}

		log.Info("Cache image doesn't saved", slog.Any("id", id))
		return img, data.Url, nil
	}
	return data.Data, data.Url, nil
}
