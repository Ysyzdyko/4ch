package ports

import (
	"1337b04rd/internal/application/dto"
	"context"
)

type RickyPort interface {
	GetAllCharacters(ctx context.Context) (*dto.CharacterResponse, error)
	GetCharacterByID(ctx context.Context, id int) (*dto.Character, error)
}
