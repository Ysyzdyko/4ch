package triple_s

import (
	"context"
	"net/url"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type TripleSConfig struct {
	Endpoint      string
	User          string
	Password      string
	UseSSL        bool
	DefaultBucket string
	URLExpiry     time.Duration
}

type TripleSClient struct {
	Client        *minio.Client
	DefaultBucket string
	URLExpiry     time.Duration
}

func DefaultConfig() TripleSConfig {
	raw := os.Getenv("MINIO_ENDPOINT")
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Host == "" {
		parsed = &url.URL{
			Scheme: "http",
			Host:   "minio:9000",
		}
	}

	return TripleSConfig{
		Endpoint:      parsed.Host,
		User:          "minioadmin",
		Password:      "minioadmin",
		UseSSL:        parsed.Scheme == "https",
		DefaultBucket: "default-bucket",
		URLExpiry:     24 * time.Hour,
	}
}

func NewTripleSClient(cfg TripleSConfig) (*TripleSClient, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.User, cfg.Password, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	time.Sleep(3 * time.Second)

	err = client.MakeBucket(ctx, cfg.DefaultBucket, minio.MakeBucketOptions{})
	if err != nil {
		exists, errExists := client.BucketExists(ctx, cfg.DefaultBucket)
		if errExists != nil || !exists {
			return nil, err // ошибка создания и бакета нет
		}
	}

	return &TripleSClient{
		Client:        client,
		DefaultBucket: cfg.DefaultBucket,
		URLExpiry:     cfg.URLExpiry,
	}, nil
}
