package triple_s

import (
	"context"
	"io"
	"net/url"
	"testing"
	"time"

	"github.com/minio/minio-go/v7"
)

type MockMinioClient struct {
	PutCalled           bool
	PutArgs             []interface{}
	PutShouldReturnErr  error
	PutShouldReturnInfo minio.UploadInfo

	RemoveCalled          bool
	RemoveArgs            []interface{}
	RemoveShouldReturnErr error

	PresignedURL        *url.URL
	PresignedShouldErr  error
	PresignedCalledWith []interface{}

	StatCalled         bool
	StatShouldErr      error
	StatShouldReturn   minio.ObjectInfo
	StatCalledWithArgs []interface{}
}

func (m *MockMinioClient) PutObject(ctx context.Context, bucket, object string, reader io.Reader, size int64, opts minio.PutObjectOptions) (minio.UploadInfo, error) {
	m.PutCalled = true
	m.PutArgs = []interface{}{bucket, object, size, opts}
	return m.PutShouldReturnInfo, m.PutShouldReturnErr
}

func (m *MockMinioClient) RemoveObject(ctx context.Context, bucket, object string, opts minio.RemoveObjectOptions) error {
	m.RemoveCalled = true
	m.RemoveArgs = []interface{}{bucket, object, opts}
	return m.RemoveShouldReturnErr
}

func (m *MockMinioClient) PresignedGetObject(ctx context.Context, bucket, object string, expiry time.Duration, params url.Values) (*url.URL, error) {
	m.PresignedCalledWith = []interface{}{bucket, object, expiry, params}
	return m.PresignedURL, m.PresignedShouldErr
}

func (m *MockMinioClient) StatObject(ctx context.Context, bucket, object string, opts minio.StatObjectOptions) (minio.ObjectInfo, error) {
	m.StatCalled = true
	m.StatCalledWithArgs = []interface{}{bucket, object, opts}
	return m.StatShouldReturn, m.StatShouldErr
}

func TestAddObject_Success(t *testing.T) {
	mock := &MockMinioClient{
		PutShouldReturnErr: nil,
		PresignedURL:       &url.URL{Scheme: "https", Host: "min.io", Path: "/bucket/file.txt"},
		PresignedShouldErr: nil,
	}
	client := &TripleSClient{
		Client:        mock,
		DefaultBucket: "bucket",
		URLExpiry:     time.Hour,
	}

	data := []byte("hello")
	result, err := client.AddObject(context.Background(), "", "file.txt", data, "text/plain")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if result != "https://min.io/bucket/file.txt" {
		t.Errorf("unexpected URL: %s", result)
	}
	if !mock.PutCalled {
		t.Error("expected PutObject to be called")
	}
}
func TestAddObject_EmptyData(t *testing.T) {
	client := &TripleSClient{}
	_, err := client.AddObject(context.Background(), "bucket", "file.txt", []byte{}, "text/plain")
	if err == nil {
		t.Error("expected error on empty data")
	}
}
func TestDeleteObject_Success(t *testing.T) {
	mock := &MockMinioClient{}
	client := &TripleSClient{
		Client:        mock,
		DefaultBucket: "bucket",
	}

	err := client.DeleteObject(context.Background(), "", "file.txt")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if !mock.RemoveCalled {
		t.Error("expected RemoveObject to be called")
	}
}
func TestObjectExists_Found(t *testing.T) {
	mock := &MockMinioClient{
		StatShouldErr: nil,
	}
	client := &TripleSClient{
		Client:        mock,
		DefaultBucket: "bucket",
	}

	ok, err := client.ObjectExists(context.Background(), "", "file.txt")
	if err != nil {
		t.Fatal("unexpected error:", err)
	}
	if !ok {
		t.Error("expected object to exist")
	}
}
