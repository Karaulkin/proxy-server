package downloder

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
)

type Downloader struct {
	log *slog.Logger
}

type YouTubeDownloader interface {
	Download(videoURL string) ([]byte, error)
}

func New(log *slog.Logger) *Downloader {
	return &Downloader{
		log: log,
	}
}

// ExtractVideoID — извлекает идентификатор видео из URL
func ExtractVideoID(videoURL string) (string, error) {
	re := regexp.MustCompile(`(?:v=|youtu\.be/|youtube\.com/embed/)([a-zA-Z0-9_-]{11})`)
	matches := re.FindStringSubmatch(videoURL)
	if len(matches) < 2 {
		return "", errors.New("failed to extract the video ID by youtube")
	}
	return matches[1], nil
}

// Download — скачивает превью для указанного видео
func (d *Downloader) Download(videoURL string) ([]byte, error) {
	// Извлечение ID видео
	videoID, err := ExtractVideoID(videoURL)
	if err != nil {
		return nil, fmt.Errorf("video ID extraction error: %w", err)
	}

	// Формирование URL для превью
	thumbnailURL := fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", videoID)

	// Скачивание превью
	resp, err := http.Get(thumbnailURL)
	if err != nil {
		return nil, fmt.Errorf("preview download error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download preview, HTTP status: %d", resp.StatusCode)
	}

	// Чтение данных в []byte
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading the response: %w", err)
	}

	return data, nil
}
