package league

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// League is a league model.
// League includes league ID (name) and a set of league rules.
type League struct {
	ID    string
	Rules []LeagueRule
}

// LeagueRule is a league rule model.
// LeagueRule includes league ID. (That's all we need actually)
type LeagueRule struct {
	ID string
}

const (
	getLeaguesApiString = "https://pathofexile.com/api/leagues?type=main&realm=pc"
	// Solo Self-Found league ID. These leagues are of no interest because there is no trade.
	ssfLeagueID = "NoParties"
)

// FetchLeagues fetches all current main (no events, no season leagues) leagues using PoE API.
// Currently only PC Realm leagues are fetched.
func FetchLeagues() ([]League, error) {
	resp, err := http.Get(getLeaguesApiString)
	if err != nil {
		return nil, fmt.Errorf("FetchTradeLeagues: http protocol error: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("FetchTradeLeagues: unhandled error while fetching leagues")
	}

	var leagues []League

	if err := json.NewDecoder(resp.Body).Decode(&leagues); err != nil {
		return nil, fmt.Errorf("FetchTradeLeagues: json unmarshalling error: %s", err.Error())
	}

	return leagues, nil
}

// GetTradeLeagues filters leagues and returns trade (non-SSF) leagues.
// Trade leagues are distinguished by not having NoParties rule.
func GetTradeLeagues(leagues []League) []League {
	tradeLeagues := []League{}
NextLeague:
	for _, league := range leagues {
		for _, rule := range league.Rules {
			if rule.ID == ssfLeagueID {
				continue NextLeague
			}
		}
		tradeLeagues = append(tradeLeagues, league)
	}

	return tradeLeagues
}
