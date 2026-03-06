package maps

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	infoFieldsBase    = "id,displayName,formattedAddress,location,rating,userRatingCount,priceLevel,types,regularOpeningHours,currentOpeningHours,nationalPhoneNumber,websiteUri"
	infoFieldsFeedback = "reviews"
	infoFieldsImages   = "photos"
)

// PlaceDetails fetches full details for a place by its ID.
func (g *GoogleMaps) PlaceDetails(ctx context.Context, in PlaceInfoInput) (PlaceInfo, error) {
	pid := trimmed(in.ID)
	if pid == "" {
		return PlaceInfo{}, InputError{Param: "id", Reason: "cannot be empty"}
	}

	ep, err := g.endpoint("/places/"+pid, map[string]string{
		"languageCode": trimmed(in.Lang),
		"regionCode":   trimmed(in.CountryCode),
	})
	if err != nil {
		return PlaceInfo{}, err
	}

	raw, err := g.call(ctx, http.MethodGet, ep, nil, infoFieldMask(in))
	if err != nil {
		return PlaceInfo{}, err
	}

	var p apiPlace
	if err := json.Unmarshal(raw, &p); err != nil {
		return PlaceInfo{}, fmt.Errorf("googlemapscli: unmarshal details: %w", err)
	}

	return toPlaceInfo(p), nil
}

func infoFieldMask(in PlaceInfoInput) string {
	parts := []string{infoFieldsBase}
	if in.WithReviews {
		parts = append(parts, infoFieldsFeedback)
	}
	if in.WithPhotos {
		parts = append(parts, infoFieldsImages)
	}
	return strings.Join(parts, ",")
}
