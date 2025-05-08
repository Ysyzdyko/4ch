package ports

import (
	"context"
)

type MinioPort interface {
	Object
}

type Object interface {
	AddObject(ctx context.Context, bucket, objectName string, data []byte, contentType string) (string, error)
	DeleteObject(ctx context.Context, bucket, objectName string) (err error)
	ObjectExists(ctx context.Context, bucket, objectName string) (bool, error)
}
