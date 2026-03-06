package maps

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const nearbyFields = "places.id,places.displayName,places.formattedAddress,places.location,places.rating,places.userRatingCount,places.priceLevel,places.types,places.currentOpeningHours"

// NearbySearch finds places around a geographic point.
func (g *GoogleMaps) NearbySearch(ctx context.Context, in NearbyInput) (NearbyOutput, error) {
	in = applyNearbyDefaults(in)
	if err := checkNearbyInput(in); err != nil {
		return NearbyOutput{}, err
	}

	body := map[string]any{
		"locationRestriction": areaToCircle(in.Center),
		"maxResultCount":      in.MaxResults,
	}
	if trimmed(in.Lang) != "" {
		body["languageCode"] = trimmed(in.Lang)
	}
	if trimmed(in.CountryCode) != "" {
		body["regionCode"] = trimmed(in.CountryCode)
	}
	if len(in.Include) > 0 {
		body["includedTypes"] = in.Include
	}
	if len(in.Exclude) > 0 {
		body["excludedTypes"] = in.Exclude
	}

	ep, err := g.endpoint("/places:searchNearby", nil)
	if err != nil {
		return NearbyOutput{}, err
	}

	raw, err := g.call(ctx, http.MethodPost, ep, body, nearbyFields)
	if err != nil {
		return NearbyOutput{}, err
	}

	var apiResp apiSearchResult
	if err := json.Unmarshal(raw, &apiResp); err != nil {
		return NearbyOutput{}, fmt.Errorf("googlemapscli: unmarshal nearby: %w", err)
	}

	places := make([]PlaceBrief, 0, len(apiResp.Places))
	for _, p := range apiResp.Places {
		places = append(places, toBrief(p))
	}

	return NearbyOutput{Places: places, NextCursor: apiResp.NextPageToken}, nil
}

func applyNearbyDefaults(in NearbyInput) NearbyInput {
	if in.MaxResults == 0 {
		in.MaxResults = NearbyDefaultResults
	}
	return in
}
