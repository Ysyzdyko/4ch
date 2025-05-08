package api

import (
	"1337b04rd/internal/adapters/right/ricky"
	"1337b04rd/internal/adapters/right/triple_s"
	"1337b04rd/internal/application/core"
	"context"
	"testing"
	"time"
)

// Мок базы данных
func TestGetSessionByID(t *testing.T) {
	app := NewApp(&triple_s.TripleSClient{}, &core.User{}, &MockDbPort{}, ricky.NewRicky())

	tests := []struct {
		sessionID  string
		wantError  bool
		wantActive bool
	}{
		{
			sessionID:  "valid-session",
			wantError:  false,
			wantActive: true, // Сессия должна быть активной, так как expiresAt в будущем
		},
		{
			sessionID:  "invalid-session",
			wantError:  true,
			wantActive: false, // Если сессии нет в базе, ошибки не будет, вернётся nil
		},
	}

	for _, tt := range tests {
		t.Run(tt.sessionID, func(t *testing.T) {
			session, err := app.GetSessionByID(context.Background(), tt.sessionID)

			// Проверяем ошибки
			if (err != nil) != tt.wantError {
				t.Errorf("GetSessionByID() error = %v, wantError %v", err, tt.wantError)
			}

			// Проверяем значение IsActive
			if session != nil && session.IsActive != tt.wantActive {
				t.Errorf("GetSessionByID() = %v, wantActive %v", session.IsActive, tt.wantActive)
			}
		})
	}
}

func TestCreateSession(t *testing.T) {
	// Инициализируем приложение с мок-базой данных
	app := NewApp(&triple_s.TripleSClient{}, &core.User{}, &MockDbPort{}, ricky.NewRicky())

	// Тестируем создание сессии
	session, err := app.CreateSession(context.Background())
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	// Проверяем, что сессия создана
	if session.ID == "" {
		t.Errorf("CreateSession() = %v, want non-empty ID", session)
	}

	// Проверяем, что сессия активна
	if !session.IsActive {
		t.Errorf("CreateSession() = %v, want IsActive true", session)
	}

	// Проверяем, что ExpiresAt установлено через 10 минут
	if session.ExpiresAt.Before(time.Now().Add(9*time.Minute)) || session.ExpiresAt.After(time.Now().Add(11*time.Minute)) {
		t.Errorf("CreateSession() = %v, want ExpiresAt to be around 10 minutes from now", session.ExpiresAt)
	}
}
