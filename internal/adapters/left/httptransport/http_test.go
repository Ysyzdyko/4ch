package http

import (
	"1337b04rd/internal/application/domain"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewHTTPServer(t *testing.T) {
	t.Parallel()

	testServer := httptest.NewServer(newRouter(&fakeAPIPort{}))

	defer testServer.Close()

	// пример запроса на health-check или корень
	resp, err := http.Get(testServer.URL + "/catalog")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 OK, got %v", resp.Status)
	}
}

type fakeAPIPort struct{}

func (f *fakeAPIPort) GetCatalog(ctx context.Context) ([]*domain.PostSummary, error) {
	return []*domain.PostSummary{
		{ID: "1", Title: "Test Post", CreatedAt: time.Now()},
	}, nil
}

func (f *fakeAPIPort) GetPostByID(ctx context.Context, id string) (*domain.Post, error) {
	return nil, nil
}
func (f *fakeAPIPort) GetArchiveList(ctx context.Context) ([]*domain.PostSummary, error) {
	return nil, nil
}
func (f *fakeAPIPort) GetArchivedPostByID(ctx context.Context, id string) (*domain.Post, error) {
	return nil, nil
}
func (f *fakeAPIPort) CreatePost(ctx context.Context, post *domain.Post) error {
	return nil
}
func (f *fakeAPIPort) AddComment(ctx context.Context, postID string, comment *domain.Comment) error {
	return nil
}
func (f *fakeAPIPort) ReplyToComment(ctx context.Context, parentCommentID string, reply *domain.Comment) error {
	return nil
}
func (f *fakeAPIPort) GetSessionByID(ctx context.Context, sessionID string) (*domain.Session, error) {
	return nil, nil
}
func (f *fakeAPIPort) CreateSession(ctx context.Context) (*domain.Session, error) {
	return nil, nil
}

func (f *fakeAPIPort) UploadAndCreatePost(ctx context.Context, post *domain.Post, objectName string, fileData []byte, contentType string) error {
	return nil
}
