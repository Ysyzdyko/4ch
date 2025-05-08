package ricky

import (
	"1337b04rd/internal/application/dto"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	ports "1337b04rd/internal/ports/right"
)

const (
	userAgent            = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36"
	baseURL              = "https://rickandmortyapi.com/api/character"
	getCharacterByIDTmpl = "https://rickandmortyapi.com/api/character/%d"
)

type RickyAndMorty struct {
	client *http.Client
}

func NewRicky() ports.RickyPort {
	return &RickyAndMorty{
		client: &http.Client{
			Timeout: 100 * time.Second,
		},
	}
}

func (r *RickyAndMorty) GetAllCharacters(ctx context.Context) (*dto.CharacterResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Important: Set the User-Agent header to avoid being blocked by the API
	req.Header.Set("User-Agent", userAgent)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	if resp.Body == nil {
		return nil, fmt.Errorf("empty response body")
	}

	var charResp dto.CharacterResponse
	if err := json.NewDecoder(resp.Body).Decode(&charResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &charResp, nil
}

func (r *RickyAndMorty) GetCharacterByID(ctx context.Context, id int) (*dto.Character, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(getCharacterByIDTmpl, id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Important: Set the User-Agent header to avoid being blocked by the API
	req.Header.Set("User-Agent", userAgent)

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}

	if resp.Body == nil {
		return nil, fmt.Errorf("empty response body")
	}

	var char dto.Character
	if err := json.NewDecoder(resp.Body).Decode(&char); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &char, nil
}
