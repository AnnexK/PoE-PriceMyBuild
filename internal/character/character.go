package character

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AnnexK/PoE-PriceMyBuild/internal/league"
)

// Character is a character model.
// Character includes character name and league name.
type Character struct {
	Name   string
	League string
}

const getCharactersApiString = "https://pathofexile.com/character-window/get-characters?accountName=%s"

// FetchCharacters fetches all characters for an account `username` using PoE API.
// If the account is private, an HTTP 403 is returned.
func FetchCharacters(username string) ([]Character, error) {
	resp, err := http.Get(fmt.Sprintf(getCharactersApiString, username))
	if err != nil {
		return nil, fmt.Errorf("FetchTradeCharacters: http protocol error: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode == 403 {
		return nil, fmt.Errorf("FetchTradeCharacters: account %s is private", username)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("FetchTradeCharacters: unhandled error while fetching characters")
	}

	var characters []Character

	if err := json.NewDecoder(resp.Body).Decode(&characters); err != nil {
		return nil, fmt.Errorf("FetchTradeCharacters: json unmarshalling error: %s", err.Error())
	}

	return characters, nil
}

// GetCharactersInLeagues filters a character slice by league name.
// It returns a slice of characters that belong to at least one league in `leagues`.
func GetCharactersInLeagues(characters []Character, leagues []league.League) []Character {
	charactersInLeagues := []Character{}

	for _, character := range characters {
		for _, league := range leagues {
			if character.League == league.ID {
				charactersInLeagues = append(charactersInLeagues, character)
				break
			}
		}
	}
	return charactersInLeagues
}
