package api

import (
	"1337b04rd/internal/application/domain"
	"1337b04rd/pkg/utils"
	"context"
	"time"
)

func (app *App) getUser(ctx context.Context, userID string) error {
	return nil
}

func (app *App) createUser(ctx context.Context) (*domain.User, error) {
	charResp, err := app.ricky.GetAllCharacters(ctx)
	if err != nil {
		return nil, err
	}

	// Get the maximum character ID from the response
	maxCharID := charResp.GetCharactersCount()

	// Get the maximum character ID from the database
	userCount, err := app.db.GetMaxCharacterID(ctx)
	if err != nil {
		return nil, err
	}

	// Generate a new character ID
	newCharID, err := app.userSvc.GenUserID(userCount, maxCharID)
	if err != nil {
		return nil, err
	}

	// Fetch the character details using the new character ID
	chr, err := app.ricky.GetCharacterByID(ctx, newCharID)
	if err != nil {
		return nil, err
	}

	userID, err := utils.GenerateUUID()
	if err != nil {
		return nil, err
	}

	newUser := &domain.User{
		ID:        userID,
		Username:  chr.GetName(),
		ImageURL:  chr.GetImageURL(),
		CreatedAt: time.Now(),
	}

	if err := app.db.CreateUser(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}
