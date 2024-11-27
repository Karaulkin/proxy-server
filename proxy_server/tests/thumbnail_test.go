package tests

import (
	"encoding/hex"
	thumbnailv1 "github.com/Karaulkin/proto_thumbnail/gen/go/thumbnail"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"proxy_server/internal/lib/downloder"
	"proxy_server/internal/lib/logger/handlers/slogpretty"
	"proxy_server/tests/suite"
	"testing"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func TestGetThumbnail_HappyPath_Download(t *testing.T) {
	ctx, st := suite.New(t)
	log := setupLogger(st.Cfg.Env)

	url := "https://www.youtube.com/watch?v=JqfbfraMjW8"

	data, err := downloder.New(log).Download(url)
	if err != nil {
		t.Fatal(err)
	}

	respData, err := st.ThumbnailClient.GetThumbnail(ctx, &thumbnailv1.ThumbnailRequest{
		VideoUrl: url,
	})

	if err != nil {
		require.Error(t, err)
		require.Contains(t, err.Error(), "rpc error: code = Internal desc = video ID extraction error: failed to download image preview")
	} else {
		require.NoError(t, err)
		assert.NotEmpty(t, respData.GetVideoUrl())

		assert.Equal(t, data, respData.GetImageData())
	}

}

func TestGetThumbnail_HappyPath_LoadCash(t *testing.T) {
	ctx, st := suite.New(t)

	url := "https://www.youtube.com/watch?v=hDY5_GqLzho"

	data, err := hex.DecodeString("FFD8FFE000104A46494600010101006000600000FFD9")
	if err != nil {
		t.Fatal(err)
	}

	respData, err := st.ThumbnailClient.GetThumbnail(ctx, &thumbnailv1.ThumbnailRequest{
		VideoUrl: url,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, respData.GetVideoUrl())

	assert.Equal(t, data, respData.GetImageData())

}

func TestGetThumbnail_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		url         string
		expectedErr string
	}{
		{
			url:         "",
			expectedErr: "rpc error: code = InvalidArgument desc = missing url",
		},
		{
			url:         "some-broken-url",
			expectedErr: "rpc error: code = Internal desc = video ID extraction error: failed to extract the video ID by youtube",
		},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			_, err := st.ThumbnailClient.GetThumbnail(ctx, &thumbnailv1.ThumbnailRequest{
				VideoUrl: tt.url,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)

		})
	}
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
