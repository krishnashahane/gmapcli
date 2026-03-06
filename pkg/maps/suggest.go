package maps

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const suggestFields = "suggestions.placePrediction.placeId,suggestions.placePrediction.place,suggestions.placePrediction.text,suggestions.placePrediction.structuredFormat,suggestions.placePrediction.types,suggestions.placePrediction.distanceMeters,suggestions.queryPrediction.text,suggestions.queryPrediction.structuredFormat"

// Suggest returns autocomplete predictions for partial input.
func (g *GoogleMaps) Suggest(ctx context.Context, in SuggestInput) (SuggestOutput, error) {
	in = applySuggestDefaults(in)
	if err := checkSuggestInput(in); err != nil {
		return SuggestOutput{}, err
	}

	body := map[string]any{
		"input": trimmed(in.Fragment),
	}
	if trimmed(in.Session) != "" {
		body["sessionToken"] = trimmed(in.Session)
	}
	if trimmed(in.Lang) != "" {
		body["languageCode"] = trimmed(in.Lang)
	}
	if trimmed(in.CountryCode) != "" {
		body["regionCode"] = trimmed(in.CountryCode)
	}
	if in.Vicinity != nil {
		body["locationBias"] = areaToCircle(in.Vicinity)
	}

	ep, err := g.endpoint("/places:autocomplete", nil)
	if err != nil {
		return SuggestOutput{}, err
	}

	raw, err := g.call(ctx, http.MethodPost, ep, body, suggestFields)
	if err != nil {
		return SuggestOutput{}, err
	}

	var apiResp apiSuggestResult
	if err := json.Unmarshal(raw, &apiResp); err != nil {
		return SuggestOutput{}, fmt.Errorf("googlemapscli: unmarshal suggest: %w", err)
	}

	items := make([]Suggestion, 0, len(apiResp.Suggestions))
	for _, s := range apiResp.Suggestions {
		if mapped, ok := convertSuggestion(s); ok {
			items = append(items, mapped)
		}
	}

	if in.MaxResults > 0 && len(items) > in.MaxResults {
		items = items[:in.MaxResults]
	}

	return SuggestOutput{Items: items}, nil
}

func convertSuggestion(s apiSuggestionItem) (Suggestion, bool) {
	if s.PlacePrediction != nil {
		pp := s.PlacePrediction
		sf := pp.StructuredFormat
		return Suggestion{
			Type:       "place",
			PlaceID:    pp.PlaceId,
			Resource:   pp.Place,
			Label:      smallText(pp.Text),
			Primary:    smallText(structuredMain(sf)),
			Secondary:  smallText(structuredSub(sf)),
			Categories: pp.Types,
			DistanceM:  pp.DistanceMeters,
		}, true
	}
	if s.QueryPrediction != nil {
		qp := s.QueryPrediction
		sf := qp.StructuredFormat
		return Suggestion{
			Type:      "query",
			Label:     smallText(qp.Text),
			Primary:   smallText(structuredMain(sf)),
			Secondary: smallText(structuredSub(sf)),
		}, true
	}
	return Suggestion{}, false
}

func structuredMain(sf *apiStructuredText) *apiSmallText {
	if sf == nil {
		return nil
	}
	return sf.MainText
}

func structuredSub(sf *apiStructuredText) *apiSmallText {
	if sf == nil {
		return nil
	}
	return sf.SecondaryText
}

func smallText(t *apiSmallText) string {
	if t == nil {
		return ""
	}
	return t.Text
}

func applySuggestDefaults(in SuggestInput) SuggestInput {
	if in.MaxResults == 0 {
		in.MaxResults = SuggestDefaultResults
	}
	return in
}
