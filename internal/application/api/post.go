package api

import (
	"1337b04rd/internal/application/domain"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx"
)

func (app *App) CreatePost(ctx context.Context, post *domain.Post) error {
	app.Lock()
	defer app.Unlock()

	if err := app.db.CreatePost(ctx, post); err != nil {
		return err
	}

	// Таймер
	if timer, ok := app.timers[post.ID]; ok {
		timer.Stop()
	}
	app.timers[post.ID] = time.AfterFunc(10*time.Minute, func() {
		app.Lock()
		defer app.Unlock()

		if timer, ok := app.timers[post.ID]; ok {
			delete(app.timers, post.ID)
			go app.archivePost(context.Background(), post.ID)
			timer.Stop()
		}
	})

	return nil
}

func (app *App) GetPostByID(ctx context.Context, id string) (*domain.Post, error) {
	post, err := app.db.GetPostByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get post by id: %w", err)
	}

	author, err := app.db.GetUserByID(ctx, post.AuthorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	post.UserAvatar = author.ImageURL

	return post, nil
}

func (app *App) GetCatalog(ctx context.Context) ([]*domain.PostSummary, error) {
	catalog, err := app.db.ListCatalog(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return catalog, nil
}

func (app *App) GetArchiveList(ctx context.Context) ([]*domain.PostSummary, error) {
	catalog, err := app.db.ListArchiveCatalog(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return catalog, nil
}

func (app *App) GetArchivedPostByID(ctx context.Context, id string) (*domain.Post, error) {
	post, err := app.db.GetArchivedPostByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}

		return nil, err
	}

	author, err := app.db.GetUserByID(ctx, post.AuthorID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	post.UserAvatar = author.ImageURL

	return post, nil
}

func (app *App) archivePost(ctx context.Context, id string) (*domain.Post, error) {
	return app.db.ArchivePostByID(ctx, id)
}

func (app *App) UploadAndCreatePost(ctx context.Context, post *domain.Post, objectName string, fileData []byte, contentType string) error {
	// Проверяем, не загружен ли файл уже
	flag, err := app.minio.ObjectExists(ctx, "", objectName)
	if err != nil {
		return err
	}
	if !flag {
		url, err := app.minio.AddObject(ctx, "", objectName, fileData, contentType)
		if err != nil {
			return err
		}

		post.ImageURL = url
	}

	// Только база и таймер
	return app.CreatePost(ctx, post)
}
