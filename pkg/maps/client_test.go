package maps

import (
	"net/http"
	"testing"
	"time"
)

func TestNormalizeBaseURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		raw string
		fallback string
		want string
	}{
		{"official places", "https://places.googleapis.com/v1/", PlacesEndpoint, "https://places.googleapis.com/v1"},
		{"official routes", "https://routes.googleapis.com/", RoutesEndpoint, "https://routes.googleapis.com"},
		{"untrusted host falls back", "https://example.com", PlacesEndpoint, PlacesEndpoint},
		{"http falls back", "http://127.0.0.1:8080", PlacesEndpoint, PlacesEndpoint},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeBaseURL(tt.raw, tt.fallback); got != tt.want {
				t.Fatalf("normalizeBaseURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNewGoogleMapsDefaultTimeout(t *testing.T) {
	g := NewGoogleMaps(Settings{})
	if g.http == nil || g.http.Timeout != 10*time.Second {
		t.Fatalf("unexpected HTTP client timeout: %#v", g.http)
	}
	if g.http.CheckRedirect == nil {
		t.Fatal("expected redirect protection")
	}
	if err := g.http.CheckRedirect(&http.Request{}, nil); err != http.ErrUseLastResponse {
		t.Fatalf("CheckRedirect() = %v, want %v", err, http.ErrUseLastResponse)
	}
}
