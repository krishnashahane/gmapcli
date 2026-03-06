package maps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const navFields = "routes.description,routes.warnings,routes.legs.distanceMeters,routes.legs.duration,routes.legs.localizedValues.distance,routes.legs.localizedValues.duration,routes.legs.steps.distanceMeters,routes.legs.steps.staticDuration,routes.legs.steps.localizedValues.distance,routes.legs.steps.localizedValues.staticDuration,routes.legs.steps.navigationInstruction.instructions,routes.legs.steps.navigationInstruction.maneuver,routes.legs.steps.travelMode"

const (
	ModeWalk    = "walking"
	ModeDrive   = "driving"
	ModeBike    = "bicycling"
	ModeTransit = "transit"
)

const (
	UnitMetric   = "metric"
	UnitImperial = "imperial"
)

var validUnits = map[string]struct{}{
	UnitMetric:   {},
	UnitImperial: {},
}

// Navigate computes directions between two locations via the Routes API.
func (g *GoogleMaps) Navigate(ctx context.Context, in NavigationInput) (NavigationOutput, error) {
	in = prepareNavInput(in)
	if err := verifyNavInput(in); err != nil {
		return NavigationOutput{}, err
	}

	body := buildNavBody(in)
	ep := navEndpoint(g.directionsURL)

	raw, err := g.call(ctx, http.MethodPost, ep, body, navFields)
	if err != nil {
		return NavigationOutput{}, err
	}

	var apiResp apiRoutesResult
	if err := json.Unmarshal(raw, &apiResp); err != nil {
		return NavigationOutput{}, fmt.Errorf("googlemapscli: unmarshal nav: %w", err)
	}
	if len(apiResp.Routes) == 0 || len(apiResp.Routes[0].Legs) == 0 {
		return NavigationOutput{}, errors.New("googlemapscli: no route found")
	}

	route := apiResp.Routes[0]
	leg := route.Legs[0]

	maneuvers := make([]NavStep, 0, len(leg.Steps))
	for _, s := range leg.Steps {
		maneuvers = append(maneuvers, NavStep{
			Direction:     trimmed(s.NavigationInstruction.Instructions),
			DistanceLabel: trimmed(s.LocalizedValues.Distance.Text),
			DistanceM:     s.DistanceMeters,
			DurationLabel: trimmed(s.LocalizedValues.StaticDuration.Text),
			DurationSec:   toDurationSec(s.StaticDuration),
			Method:        trimmed(s.TravelMode),
			Action:        trimmed(s.NavigationInstruction.Maneuver),
		})
	}

	return NavigationOutput{
		TravelBy:      strings.ToUpper(in.TravelBy),
		Description:   trimmed(route.Description),
		FromLabel:     navLabel(in.OriginID, in.OriginCoords, in.Origin),
		ToLabel:       navLabel(in.DestinationID, in.DestinationCoords, in.Destination),
		DistanceLabel: trimmed(leg.LocalizedValues.Distance.Text),
		DistanceM:     leg.DistanceMeters,
		DurationLabel: trimmed(leg.LocalizedValues.Duration.Text),
		DurationSec:   toDurationSec(leg.Duration),
		Alerts:        route.Warnings,
		Maneuvers:     maneuvers,
	}, nil
}

func prepareNavInput(in NavigationInput) NavigationInput {
	in.Origin = trimmed(in.Origin)
	in.Destination = trimmed(in.Destination)
	in.OriginID = trimmed(in.OriginID)
	in.DestinationID = trimmed(in.DestinationID)
	in.TravelBy = strings.ToLower(trimmed(in.TravelBy))
	if in.TravelBy == "" {
		in.TravelBy = ModeWalk
	}
	if norm := canonicalMode(in.TravelBy); norm != "" {
		in.TravelBy = norm
	}
	in.MeasureSystem = strings.ToLower(trimmed(in.MeasureSystem))
	if in.MeasureSystem == "" {
		in.MeasureSystem = UnitMetric
	}
	return in
}

func verifyNavInput(in NavigationInput) error {
	if canonicalMode(in.TravelBy) == "" {
		return InputError{Param: "travel_by", Reason: "must be walk, drive, bicycle, or transit"}
	}
	if err := verifyNavPoint("origin", in.OriginID, in.OriginCoords, in.Origin); err != nil {
		return err
	}
	if err := verifyNavPoint("destination", in.DestinationID, in.DestinationCoords, in.Destination); err != nil {
		return err
	}
	if in.MeasureSystem != "" {
		if _, ok := validUnits[in.MeasureSystem]; !ok {
			return InputError{Param: "measure_system", Reason: "must be metric or imperial"}
		}
	}
	return nil
}

func verifyNavPoint(label string, placeID string, coords *Coordinates, text string) error {
	count := 0
	if trimmed(placeID) != "" {
		count++
	}
	if coords != nil {
		count++
		if coords.Latitude < -90 || coords.Latitude > 90 {
			return InputError{Param: label + ".latitude", Reason: "must be -90..90"}
		}
		if coords.Longitude < -180 || coords.Longitude > 180 {
			return InputError{Param: label + ".longitude", Reason: "must be -180..180"}
		}
	}
	if trimmed(text) != "" {
		count++
	}
	if count == 0 {
		return InputError{Param: label, Reason: "required"}
	}
	if count > 1 {
		return InputError{Param: label, Reason: "provide only one of text, place_id, or coordinates"}
	}
	return nil
}

// CanonicalMode normalizes a travel mode string.
func CanonicalMode(mode string) string {
	return canonicalMode(mode)
}

func canonicalMode(mode string) string {
	switch strings.ToLower(trimmed(mode)) {
	case "walk", "walking":
		return ModeWalk
	case "drive", "driving":
		return ModeDrive
	case "bike", "bicycle", "bicycling":
		return ModeBike
	case "transit":
		return ModeTransit
	default:
		return ""
	}
}

func navLabel(placeID string, coords *Coordinates, text string) string {
	if trimmed(placeID) != "" {
		return "place_id:" + trimmed(placeID)
	}
	if coords != nil {
		return fmt.Sprintf("%.6f,%.6f", coords.Latitude, coords.Longitude)
	}
	return trimmed(text)
}

const routeComputePath = "/directions/v2:computeRoutes"

func navEndpoint(base string) string {
	if strings.HasSuffix(base, routeComputePath) {
		return base
	}
	return base + routeComputePath
}

func modeToAPITravel(mode string) string {
	switch canonicalMode(mode) {
	case ModeWalk:
		return "WALK"
	case ModeDrive:
		return "DRIVE"
	case ModeBike:
		return "BICYCLE"
	case ModeTransit:
		return "TRANSIT"
	default:
		return "WALK"
	}
}

func measureToAPI(measure string) string {
	if strings.ToLower(trimmed(measure)) == UnitImperial {
		return "IMPERIAL"
	}
	return "METRIC"
}

func navWaypoint(placeID string, coords *Coordinates, text string) map[string]any {
	if id := trimmed(placeID); id != "" {
		return map[string]any{"placeId": id}
	}
	if coords != nil {
		return map[string]any{
			"location": map[string]any{
				"latLng": map[string]any{
					"latitude":  coords.Latitude,
					"longitude": coords.Longitude,
				},
			},
		}
	}
	return map[string]any{"address": trimmed(text)}
}

func buildNavBody(in NavigationInput) map[string]any {
	body := map[string]any{
		"origin":      navWaypoint(in.OriginID, in.OriginCoords, in.Origin),
		"destination": navWaypoint(in.DestinationID, in.DestinationCoords, in.Destination),
		"travelMode":  modeToAPITravel(in.TravelBy),
		"units":       measureToAPI(in.MeasureSystem),
	}
	if trimmed(in.Lang) != "" {
		body["languageCode"] = trimmed(in.Lang)
	}
	if trimmed(in.CountryCode) != "" {
		body["regionCode"] = trimmed(in.CountryCode)
	}
	return body
}

func toDurationSec(raw string) int {
	d, err := time.ParseDuration(trimmed(raw))
	if err != nil {
		return 0
	}
	return int(d.Seconds())
}
