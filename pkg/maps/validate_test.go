package maps

import "testing"

func TestCheckArea(t *testing.T) {
	t.Parallel()

	valid := &Area{Latitude: 19.076, Longitude: 72.8777, Radius: 100}
	if err := checkArea(valid); err != nil {
		t.Fatalf("checkArea(valid) = %v", err)
	}

	invalid := []*Area{
		{Latitude: 91, Longitude: 0, Radius: 100},
		{Latitude: 0, Longitude: 181, Radius: 100},
		{Latitude: 0, Longitude: 0, Radius: 0},
	}
	for _, area := range invalid {
		if err := checkArea(area); err == nil {
			t.Fatalf("checkArea(%+v) expected an error", area)
		}
	}
}

func TestNavigationPointValidation(t *testing.T) {
	t.Parallel()

	if err := verifyNavPoint("origin", "", nil, "Mumbai"); err != nil {
		t.Fatalf("text origin rejected: %v", err)
	}
	if err := verifyNavPoint("origin", "ChIJ123", nil, ""); err != nil {
		t.Fatalf("place ID origin rejected: %v", err)
	}
	if err := verifyNavPoint("origin", "", &Coordinates{Latitude: 19, Longitude: 73}, ""); err != nil {
		t.Fatalf("coordinate origin rejected: %v", err)
	}
	if err := verifyNavPoint("origin", "ChIJ123", nil, "Mumbai"); err == nil {
		t.Fatal("expected mutually exclusive origin inputs to fail")
	}
	if err := verifyNavPoint("origin", "", nil, ""); err == nil {
		t.Fatal("expected missing origin to fail")
	}
}
