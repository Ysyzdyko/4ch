package ricky

import (
	"context"
	"testing"
	"time"
)

func TestRickyAndMorty_GetAllCharacters(t *testing.T) {
	r := NewRicky()

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	chrResp, err := r.GetAllCharacters(ctx)
	if err != nil {
		t.Errorf("failed to get max character ID: %v", err)
	}

	t.Logf("max character ID: %v", chrResp.GetCharactersCount())
}

func TestRickyAndMorty_GetCharacterByID(t *testing.T) {
	r := NewRicky()

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Second)
	defer cancel()

	id := 1

	chr, err := r.GetCharacterByID(ctx, id)
	if err != nil {
		t.Errorf("failed to get character by ID: %v", err)
	}

	t.Logf("character: %v", chr.GetName())
}
