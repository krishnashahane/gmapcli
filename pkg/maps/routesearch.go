package maps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"sort"
	"strings"
)

const polylineFields = "routes.polyline.encodedPolyline"

const (
	earthRadius     = 6371000.0
	polylineFactor  = 1e5
)

var apiTravelModes = map[string]struct{}{
	"DRIVE": {}, "WALK": {}, "BICYCLE": {}, "TWO_WHEELER": {}, "TRANSIT": {},
}

// RouteSearch finds places along a driving/walking path between two points.
func (g *GoogleMaps) RouteSearch(ctx context.Context, in RouteSearchInput) (RouteSearchOutput, error) {
	in = prepareRouteInput(in)
	if err := verifyRouteInput(in); err != nil {
		return RouteSearchOutput{}, err
	}

	encoded, err := g.fetchPolyline(ctx, in)
	if err != nil {
		return RouteSearchOutput{}, err
	}

	decoded, err := parsePolyline(encoded)
	if err != nil {
		return RouteSearchOutput{}, err
	}

	samples := pickSamples(decoded, in.Stops)
	if len(samples) == 0 {
		return RouteSearchOutput{}, errors.New("googlemapscli: no route samples")
	}

	stops := make([]RouteStop, 0, len(samples))
	for _, pt := range samples {
		sr, err := g.TextSearch(ctx, TextSearchInput{
			Text:       in.Text,
			MaxResults: in.PerStop,
			Lang:       in.Lang,
			CountryCode: in.CountryCode,
			Vicinity: &Area{
				Latitude:  pt.Latitude,
				Longitude: pt.Longitude,
				Radius:    in.SearchRadius,
			},
		})
		if err != nil {
			return RouteSearchOutput{}, err
		}
		stops = append(stops, RouteStop{Point: pt, Places: sr.Places})
	}

	return RouteSearchOutput{Stops: stops}, nil
}

func prepareRouteInput(in RouteSearchInput) RouteSearchInput {
	in.Text = trimmed(in.Text)
	in.StartPoint = trimmed(in.StartPoint)
	in.EndPoint = trimmed(in.EndPoint)
	in.TravelBy = strings.ToUpper(trimmed(in.TravelBy))
	if in.TravelBy == "" {
		in.TravelBy = "DRIVE"
	}
	if in.PerStop == 0 {
		in.PerStop = RouteDefaultPerStop
	}
	if in.SearchRadius == 0 {
		in.SearchRadius = RouteDefaultSearchArea
	}
	if in.Stops == 0 {
		in.Stops = RouteDefaultSamples
	}
	return in
}

func verifyRouteInput(in RouteSearchInput) error {
	if in.Text == "" {
		return InputError{Param: "text", Reason: "cannot be empty"}
	}
	if in.StartPoint == "" {
		return InputError{Param: "start_point", Reason: "cannot be empty"}
	}
	if in.EndPoint == "" {
		return InputError{Param: "end_point", Reason: "cannot be empty"}
	}
	if in.PerStop < 1 || in.PerStop > SearchMaxResults {
		return InputError{Param: "per_stop", Reason: fmt.Sprintf("must be 1-%d", SearchMaxResults)}
	}
	if in.SearchRadius <= 0 {
		return InputError{Param: "search_radius", Reason: "must be positive"}
	}
	if in.Stops < 1 || in.Stops > RouteMaxSamples {
		return InputError{Param: "stops", Reason: fmt.Sprintf("must be 1-%d", RouteMaxSamples)}
	}
	if _, ok := apiTravelModes[in.TravelBy]; !ok {
		return InputError{Param: "travel_by", Reason: "must be DRIVE, WALK, BICYCLE, TWO_WHEELER, or TRANSIT"}
	}
	return nil
}

func (g *GoogleMaps) fetchPolyline(ctx context.Context, in RouteSearchInput) (string, error) {
	body := map[string]any{
		"origin":           map[string]any{"address": in.StartPoint},
		"destination":      map[string]any{"address": in.EndPoint},
		"travelMode":       in.TravelBy,
		"polylineQuality":  "OVERVIEW",
		"polylineEncoding": "ENCODED_POLYLINE",
	}
	if in.Lang != "" {
		body["languageCode"] = in.Lang
	}
	if in.CountryCode != "" {
		body["regionCode"] = in.CountryCode
	}

	ep := g.routesURL + routeComputePath
	raw, err := g.call(ctx, http.MethodPost, ep, body, polylineFields)
	if err != nil {
		return "", err
	}

	var resp apiRoutesResult
	if err := json.Unmarshal(raw, &resp); err != nil {
		return "", fmt.Errorf("googlemapscli: unmarshal polyline: %w", err)
	}
	if len(resp.Routes) == 0 {
		return "", errors.New("googlemapscli: no route returned")
	}
	line := trimmed(resp.Routes[0].Polyline.EncodedPolyline)
	if line == "" {
		return "", errors.New("googlemapscli: empty polyline")
	}
	return line, nil
}

func parsePolyline(encoded string) ([]Coordinates, error) {
	if trimmed(encoded) == "" {
		return nil, errors.New("googlemapscli: blank polyline")
	}
	pts := make([]Coordinates, 0, len(encoded)/4)
	var lat, lng int
	for idx := 0; idx < len(encoded); {
		var d int
		var shift uint
		for {
			if idx >= len(encoded) {
				return nil, errors.New("googlemapscli: truncated polyline")
			}
			b := int(encoded[idx]) - 63
			idx++
			d |= (b & 0x1f) << shift
			shift += 5
			if b < 0x20 {
				break
			}
		}
		lat += (d >> 1) ^ (-(d & 1))

		d = 0
		shift = 0
		for {
			if idx >= len(encoded) {
				return nil, errors.New("googlemapscli: truncated polyline")
			}
			b := int(encoded[idx]) - 63
			idx++
			d |= (b & 0x1f) << shift
			shift += 5
			if b < 0x20 {
				break
			}
		}
		lng += (d >> 1) ^ (-(d & 1))

		pts = append(pts, Coordinates{
			Latitude:  float64(lat) / polylineFactor,
			Longitude: float64(lng) / polylineFactor,
		})
	}
	return pts, nil
}

func pickSamples(pts []Coordinates, count int) []Coordinates {
	if len(pts) == 0 || count <= 0 {
		return nil
	}
	if len(pts) == 1 {
		return []Coordinates{pts[0]}
	}
	if count == 1 {
		return []Coordinates{interpolateAt(pts, pathLength(pts)/2)}
	}
	if count >= len(pts) {
		return dedup(pts)
	}

	dists := runningDistances(pts)
	total := dists[len(dists)-1]
	if total == 0 {
		return []Coordinates{pts[0]}
	}
	gap := total / float64(count-1)

	out := make([]Coordinates, 0, count)
	for i := 0; i < count; i++ {
		target := gap * float64(i)
		pt := interpolateAtCum(pts, dists, target)
		if len(out) == 0 || !closeEnough(out[len(out)-1], pt) {
			out = append(out, pt)
		}
	}
	return out
}

func runningDistances(pts []Coordinates) []float64 {
	d := make([]float64, len(pts))
	for i := 1; i < len(pts); i++ {
		d[i] = d[i-1] + haversine(pts[i-1], pts[i])
	}
	return d
}

func pathLength(pts []Coordinates) float64 {
	if len(pts) < 2 {
		return 0
	}
	var sum float64
	for i := 1; i < len(pts); i++ {
		sum += haversine(pts[i-1], pts[i])
	}
	return sum
}

func interpolateAt(pts []Coordinates, target float64) Coordinates {
	if len(pts) == 0 {
		return Coordinates{}
	}
	return interpolateAtCum(pts, runningDistances(pts), target)
}

func interpolateAtCum(pts []Coordinates, dists []float64, target float64) Coordinates {
	if target <= 0 {
		return pts[0]
	}
	total := dists[len(dists)-1]
	if target >= total {
		return pts[len(pts)-1]
	}
	idx := sort.Search(len(dists), func(i int) bool { return dists[i] >= target })
	if idx == 0 {
		return pts[0]
	}
	prev := pts[idx-1]
	next := pts[idx]
	seg := dists[idx] - dists[idx-1]
	if seg <= 0 {
		return next
	}
	frac := (target - dists[idx-1]) / seg
	return Coordinates{
		Latitude:  prev.Latitude + (next.Latitude-prev.Latitude)*frac,
		Longitude: prev.Longitude + (next.Longitude-prev.Longitude)*frac,
	}
}

func dedup(pts []Coordinates) []Coordinates {
	out := make([]Coordinates, 0, len(pts))
	for _, p := range pts {
		if len(out) == 0 || !closeEnough(out[len(out)-1], p) {
			out = append(out, p)
		}
	}
	return out
}

func closeEnough(a, b Coordinates) bool {
	const eps = 1e-6
	return math.Abs(a.Latitude-b.Latitude) < eps && math.Abs(a.Longitude-b.Longitude) < eps
}

func haversine(a, b Coordinates) float64 {
	lat1 := a.Latitude * math.Pi / 180
	lat2 := b.Latitude * math.Pi / 180
	dLat := (b.Latitude - a.Latitude) * math.Pi / 180
	dLng := (b.Longitude - a.Longitude) * math.Pi / 180

	sLat := math.Sin(dLat / 2)
	sLng := math.Sin(dLng / 2)
	h := sLat*sLat + math.Cos(lat1)*math.Cos(lat2)*sLng*sLng
	return 2 * earthRadius * math.Asin(math.Sqrt(h))
}
