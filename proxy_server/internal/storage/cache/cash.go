package cache

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"proxy_server/internal/domain/models"
	"proxy_server/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (c *Storage) Stop() error {
	_, err := c.db.Exec("DELETE FROM thumbnails")
	if err != nil {
		return fmt.Errorf("failed to clear table: %w", err)
	}

	if err := c.db.Close(); err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}
	return nil
}

func (c *Storage) CacheUrlProvider(ctx context.Context, url string) (models.CachedData, error) {
	const op = "storage.sqlite.CacheUrlProvider"

	stmt, err := c.db.PrepareContext(ctx, "SELECT id, data, url FROM thumbnails WHERE url = ?")
	if err != nil {
		return models.CachedData{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, url)

	var cache models.CachedData
	err = row.Scan(&cache.ID, &cache.Data, &cache.Url)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.CachedData{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return models.CachedData{}, fmt.Errorf("%s: %w", op, err)
	}

	return cache, nil
}

func (c *Storage) SaveCache(ctx context.Context, url string, data []byte) (int, error) {
	const op = "storage.sqlite.HashSaver"

	stmt, err := c.db.PrepareContext(ctx, "INSERT INTO thumbnails (url, data) VALUES (?, ?) ON CONFLICT(url) DO UPDATE SET  data = excluded.data")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, url, data)
	if err != nil {
		// TODO: sqlite3.Error прописать ошибку связанную с базой данных
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int(id), nil
}
