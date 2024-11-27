package thumbnail

import (
	"context"
	thumbnailv1 "github.com/Karaulkin/proto_thumbnail/gen/go/thumbnail"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Thumbnail interface {
	Url(ctx context.Context,
		videoUrl string,
	) ([]byte, string, error)
}

type serverAPI struct {
	thumbnailv1.UnimplementedThumbnailServer
	thumbnail Thumbnail
}

func Register(gRPCServer *grpc.Server, thumbnail Thumbnail) {
	thumbnailv1.RegisterThumbnailServer(gRPCServer, &serverAPI{thumbnail: thumbnail})
}

const (
	EmptyValue  = 0
	EmptyString = ""
)

func (s *serverAPI) GetThumbnail(
	ctx context.Context,
	req *thumbnailv1.ThumbnailRequest,
) (*thumbnailv1.ThumbnailResponse, error) { // хэндлер из прото файла
	if req.GetVideoUrl() == EmptyString {
		return nil, status.Error(codes.InvalidArgument, "missing url")
	}

	// Сервисный слой Url и интерфейс к нему Thumbnail
	data, url, err := s.thumbnail.Url(ctx, req.GetVideoUrl())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &thumbnailv1.ThumbnailResponse{
		ImageData: data,
		VideoUrl:  url,
	}, nil
}
