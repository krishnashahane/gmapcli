package maps

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const searchFields = "places.id,places.displayName,places.formattedAddress,places.location,places.rating,places.userRatingCount,places.priceLevel,places.types,places.currentOpeningHours,nextPageToken"

// TextSearch finds places matching a text query.
func (g *GoogleMaps) TextSearch(ctx context.Context, in TextSearchInput) (TextSearchOutput, error) {
	in = applySearchDefaults(in)
	if err := checkSearchInput(in); err != nil {
		return TextSearchOutput{}, err
	}

	body := assembleSearchBody(in)
	ep, err := g.endpoint("/places:searchText", nil)
	if err != nil {
		return TextSearchOutput{}, err
	}

	raw, err := g.call(ctx, http.MethodPost, ep, body, searchFields)
	if err != nil {
		return TextSearchOutput{}, err
	}

	var apiResp apiSearchResult
	if err := json.Unmarshal(raw, &apiResp); err != nil {
		return TextSearchOutput{}, fmt.Errorf("googlemapscli: unmarshal search: %w", err)
	}

	places := make([]PlaceBrief, 0, len(apiResp.Places))
	for _, p := range apiResp.Places {
		places = append(places, toBrief(p))
	}

	return TextSearchOutput{Places: places, NextCursor: apiResp.NextPageToken}, nil
}

func assembleSearchBody(in TextSearchInput) map[string]any {
	query := in.Text
	if in.Filters != nil && trimmed(in.Filters.Keyword) != "" {
		query = trimmed(query + " " + in.Filters.Keyword)
	}

	body := map[string]any{
		"textQuery": query,
		"pageSize":  in.MaxResults,
	}

	if trimmed(in.Lang) != "" {
		body["languageCode"] = trimmed(in.Lang)
	}
	if trimmed(in.CountryCode) != "" {
		body["regionCode"] = trimmed(in.CountryCode)
	}
	if in.Cursor != "" {
		body["pageToken"] = in.Cursor
	}
	if in.Vicinity != nil {
		body["locationBias"] = areaToCircle(in.Vicinity)
	}

	if in.Filters != nil {
		f := in.Filters
		if len(f.Categories) > 0 {
			body["includedType"] = f.Categories[0]
		}
		if f.OnlyOpen != nil {
			body["openNow"] = *f.OnlyOpen
		}
		if f.MinScore != nil {
			body["minRating"] = *f.MinScore
		}
		if len(f.PriceTiers) > 0 {
			enums := make([]string, 0, len(f.PriceTiers))
			for _, tier := range f.PriceTiers {
				if mapped, ok := numericToPrice[tier]; ok {
					enums = append(enums, mapped)
				}
			}
			if len(enums) > 0 {
				body["priceLevels"] = enums
			}
		}
	}

	return body
}

func applySearchDefaults(in TextSearchInput) TextSearchInput {
	if in.MaxResults == 0 {
		in.MaxResults = SearchDefaultResults
	}
	return in
}
