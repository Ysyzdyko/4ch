package triple_s

import (
	"bytes"
	"context"
	"errors"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
)

// AddObject uploads data to MinIO and returns a presigned URL for public access
func (m *TripleSClient) AddObject(ctx context.Context, bucket, objectName string, data []byte, contentType string) (string, error) {
	// Validate input
	if len(data) == 0 {
		return "", errors.New("empty data provided")
	}

	if bucket == "" {
		if m.DefaultBucket == "" {
			return "", errors.New("no bucket specified")
		}
		bucket = m.DefaultBucket
	}

	// Upload to MinIO
	reader := bytes.NewReader(data)
	_, err := m.Client.PutObject(
		ctx,
		bucket,
		objectName,
		reader,
		int64(len(data)),
		minio.PutObjectOptions{
			ContentType: contentType,
		},
	)
	if err != nil {
		return "", err
	}

	// Return presigned URL
	return m.GetObjectURL(ctx, bucket, objectName)
}

// GetObjectURL generates presigned URL for accessing the object
func (m *TripleSClient) GetObjectURL(ctx context.Context, bucket, objectName string) (string, error) {
	if bucket == "" {
		if m.DefaultBucket == "" {
			return "", errors.New("no bucket specified")
		}
		bucket = m.DefaultBucket
	}

	if m.URLExpiry == 0 {
		m.URLExpiry = 7 * 24 * time.Hour // Default expiry: 1 week
	}

	reqParams := make(url.Values)
	presignedURL, err := m.Client.PresignedGetObject(
		ctx,
		bucket,
		objectName,
		m.URLExpiry,
		reqParams,
	)

	if err != nil {
		return "", err
	}

	return presignedURL.String(), nil
}

// DeleteObject removes object from storage
func (m *TripleSClient) DeleteObject(ctx context.Context, bucket, objectName string) error {
	if bucket == "" {
		if m.DefaultBucket == "" {
			return errors.New("no bucket specified")
		}
		bucket = m.DefaultBucket
	}

	err := m.Client.RemoveObject(
		ctx,
		bucket,
		objectName,
		minio.RemoveObjectOptions{
			GovernanceBypass: true, // Needed if versioning is enabled
		},
	)
	if err != nil {
		return err
	}

	return nil
}

// GetObjectURL generates presigned URL for accessing the object
// func (m *TripleSClient) GetObjectURL(ctx context.Context, bucket, objectName string) (string, error) {
// 	if bucket == "" {
// 		if m.DefaultBucket == "" {
// 			return "", errors.New("no bucket specified")
// 		}
// 		bucket = m.DefaultBucket
// 	}

// 	if m.URLExpiry == 0 {
// 		m.URLExpiry = 7 * 24 * time.Hour // Default expiry: 1 week
// 	}

// 	reqParams := make(url.Values)
// 	_, err := url.Parse(m.Client.EndpointURL().String())
// 	if err != nil {
// 		return "", err
// 	}

// 	presignedURL, err := m.Client.PresignedGetObject(
// 		ctx,
// 		bucket,
// 		objectName,
// 		m.URLExpiry,
// 		reqParams,
// 	)
// 	if err != nil {
// 		return "", err
// 	}

// 	return presignedURL.String(), nil
// }

// ObjectExists checks if object exists
func (m *TripleSClient) ObjectExists(ctx context.Context, bucket, objectName string) (bool, error) {
	if bucket == "" {
		if m.DefaultBucket == "" {
			return false, errors.New("no bucket specified")
		}
		bucket = m.DefaultBucket
	}

	_, err := m.Client.StatObject(
		ctx,
		bucket,
		objectName,
		minio.StatObjectOptions{},
	)

	if err != nil {
		if minio.ToErrorResponse(err).Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// GetObject retrieves object data
// func (m *TripleSClient) GetObject(ctx context.Context, bucket, objectName string) ([]byte, error) {
// 	if bucket == "" {
// 		if m.DefaultBucket == "" {
// 			return nil, errors.New("no bucket specified")
// 		}
// 		bucket = m.DefaultBucket
// 	}

// 	obj, err := m.Client.GetObject(
// 		ctx,
// 		bucket,
// 		objectName,
// 		minio.GetObjectOptions{},
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer obj.Close()

// 	buf := new(bytes.Buffer)
// 	if _, err := buf.ReadFrom(obj); err != nil {
// 		return nil, err
// 	}

// 	return buf.Bytes(), nil
// }
