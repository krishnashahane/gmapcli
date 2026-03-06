package maps

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const lookupFields = "places.id,places.displayName,places.formattedAddress,places.location,places.types"

// LocationLookup resolves a freeform address into place candidates.
func (g *GoogleMaps) LocationLookup(ctx context.Context, in LocationLookupInput) (LocationLookupOutput, error) {
	in = applyLookupDefaults(in)
	if err := checkLookupInput(in); err != nil {
		return LocationLookupOutput{}, err
	}

	body := map[string]any{
		"textQuery": in.Address,
		"pageSize":  in.MaxResults,
	}
	if trimmed(in.Lang) != "" {
		body["languageCode"] = trimmed(in.Lang)
	}
	if trimmed(in.CountryCode) != "" {
		body["regionCode"] = trimmed(in.CountryCode)
	}

	ep, err := g.endpoint("/places:searchText", nil)
	if err != nil {
		return LocationLookupOutput{}, err
	}

	raw, err := g.call(ctx, http.MethodPost, ep, body, lookupFields)
	if err != nil {
		return LocationLookupOutput{}, err
	}

	var apiResp apiSearchResult
	if err := json.Unmarshal(raw, &apiResp); err != nil {
		return LocationLookupOutput{}, fmt.Errorf("googlemapscli: unmarshal lookup: %w", err)
	}

	matches := make([]LocationMatch, 0, len(apiResp.Places))
	for _, p := range apiResp.Places {
		matches = append(matches, toLocationMatch(p))
	}

	return LocationLookupOutput{Matches: matches}, nil
}

func applyLookupDefaults(in LocationLookupInput) LocationLookupInput {
	if in.MaxResults == 0 {
		in.MaxResults = ResolveDefaultResults
	}
	return in
}