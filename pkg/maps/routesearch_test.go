package maps

import "testing"

func TestParsePolyline(t *testing.T) {
	t.Parallel()

	points, err := parsePolyline("_p~iF~ps|U_ulLnnqC_mqNvxq`")
	if err != nil {
		t.Fatalf("parsePolyline() error = %v", err)
	}
	if len(points) != 3 {
		t.Fatalf("parsePolyline() returned %d points, want 3", len(points))
	}

	want := []Coordinates{
		{Latitude: 38.5, Longitude: -120.2},
		{Latitude: 40.7, Longitude: -120.95},
		{Latitude: 43.252, Longitude: -126.453},
	}
	for i := range want {
		if points[i].Latitude != want[i].Latitude || points[i].Longitude != want[i].Longitude {
			t.Fatalf("point %d = %+v, want %+v", i, points[i], want[i])
		}
	}
}

func TestCanonicalMode(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"walk": ModeWalk, "walking": ModeWalk,
		"drive": ModeDrive, "driving": ModeDrive,
		"bike": ModeBike, "bicycle": ModeBike, "bicycling": ModeBike,
		"transit": ModeTransit, "unknown": "",
	}
	for input, want := range tests {
		if got := CanonicalMode(input); got != want {
			t.Fatalf("CanonicalMode(%q) = %q, want %q", input, got, want)
		}
	}
}
