package api

import (
	"1337b04rd/internal/application/domain"
	"1337b04rd/pkg/utils"
	"context"
	"time"
)

func (app *App) GetSessionByID(ctx context.Context, sessionID string) (*domain.Session, error) {
	session, err := app.db.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	session.IsActive = session.ExpiresAt.After(time.Now())
	return session, nil
}

func (app *App) CreateSession(ctx context.Context) (*domain.Session, error) {
	userCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	user, err := app.createUser(userCtx)
	if err != nil {
		return nil, err
	}

	sessionID, err := utils.GenerateUUID()
	if err != nil {
		return nil, err
	}

	session := &domain.Session{
		ID:        sessionID,
		UserID:    user.ID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(10 * time.Minute),
		IsActive:  true,
	}

	err = app.db.SaveSession(ctx, session)
	if err != nil {
		return nil, err
	}

	return session, nil
}
