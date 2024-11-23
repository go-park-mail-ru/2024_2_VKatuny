package microservicesinterface

import (
	"context"

	compressmicroservice "github.com/go-park-mail-ru/2024_2_VKatuny/microservices/compress/generated"
	"google.golang.org/grpc"
)

type ICompress interface {
	CompressAndSaveFile(ctx context.Context, in *compressmicroservice.CompressAndSaveFileInput, opts ...grpc.CallOption) (*compressmicroservice.Nothing, error)
	DeleteFile(ctx context.Context, in *compressmicroservice.DeleteFileInput, opts ...grpc.CallOption) (*compressmicroservice.Nothing, error)
}
