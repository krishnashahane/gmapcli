package maps

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// GoogleMaps provides access to the Google Maps Platform APIs.
type GoogleMaps struct {
	key           string
	placesURL     string
	routesURL     string
	directionsURL string
	http          *http.Client
}

// Settings configures a GoogleMaps instance.
type Settings struct {
	Key           string
	PlacesURL     string
	RoutesURL     string
	DirectionsURL string
	HTTP          *http.Client
	Timeout       time.Duration
}

// NewGoogleMaps creates a configured API wrapper.
func NewGoogleMaps(cfg Settings) *GoogleMaps {
	placesURL := strings.TrimRight(cfg.PlacesURL, "/")
	if placesURL == "" {
		placesURL = PlacesEndpoint
	}
	routesURL := strings.TrimRight(cfg.RoutesURL, "/")
	if routesURL == "" {
		routesURL = RoutesEndpoint
	}
	directionsURL := strings.TrimRight(cfg.DirectionsURL, "/")
	if directionsURL == "" {
		directionsURL = DirectionsEndpoint
	}

	httpClient := cfg.HTTP
	if httpClient == nil {
		timeout := cfg.Timeout
		if timeout == 0 {
			timeout = 10 * time.Second
		}
		httpClient = &http.Client{Timeout: timeout}
	}

	return &GoogleMaps{
		key:           cfg.Key,
		placesURL:     placesURL,
		routesURL:     routesURL,
		directionsURL: directionsURL,
		http:          httpClient,
	}
}

func (g *GoogleMaps) call(
	ctx context.Context,
	verb string,
	endpoint string,
	payload any,
	fields string,
) ([]byte, error) {
	if trimmed(g.key) == "" {
		return nil, ErrNoAPIKey
	}

	var body io.Reader
	if payload != nil {
		encoded, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("googlemapscli: marshal: %w", err)
		}
		body = bytes.NewReader(encoded)
	}

	req, err := http.NewRequestWithContext(ctx, verb, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("googlemapscli: build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", g.key)
	if trimmed(fields) != "" {
		req.Header.Set("X-Goog-FieldMask", fields)
	}

	resp, err := g.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("googlemapscli: http: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, fmt.Errorf("googlemapscli: read body: %w", err)
	}

	if resp.StatusCode >= http.StatusBadRequest {
		return nil, &RemoteError{Code: resp.StatusCode, Payload: trimmed(string(data))}
	}

	if len(data) == 0 {
		return nil, errors.New("googlemapscli: empty response body")
	}

	return data, nil
}

func (g *GoogleMaps) endpoint(path string, params map[string]string) (string, error) {
	full := g.placesURL + path
	if len(params) == 0 {
		return full, nil
	}

	u, err := url.Parse(full)
	if err != nil {
		return "", fmt.Errorf("googlemapscli: bad url: %w", err)
	}

	q := u.Query()
	for k, v := range params {
		if trimmed(v) == "" {
			continue
		}
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func trimmed(s string) string {
	return strings.TrimSpace(s)
}
