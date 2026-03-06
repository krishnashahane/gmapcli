package maps

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// FetchPhotoURL retrieves a direct URL for a place photo.
func (g *GoogleMaps) FetchPhotoURL(ctx context.Context, in PhotoURLInput) (PhotoURLOutput, error) {
	name := trimmed(in.ResourceName)
	if name == "" {
		return PhotoURLOutput{}, InputError{Param: "resource_name", Reason: "cannot be empty"}
	}

	path := "/" + strings.TrimPrefix(name, "/") + "/media"
	params := map[string]string{"skipHttpRedirect": "true"}
	if in.MaxWidth > 0 {
		params["maxWidthPx"] = strconv.Itoa(in.MaxWidth)
	}
	if in.MaxHeight > 0 {
		params["maxHeightPx"] = strconv.Itoa(in.MaxHeight)
	}

	ep, err := g.endpoint(path, params)
	if err != nil {
		return PhotoURLOutput{}, err
	}

	raw, err := g.call(ctx, http.MethodGet, ep, nil, "")
	if err != nil {
		return PhotoURLOutput{}, err
	}

	var resp apiPhotoMedia
	if err := json.Unmarshal(raw, &resp); err != nil {
		return PhotoURLOutput{}, fmt.Errorf("googlemapscli: unmarshal photo: %w", err)
	}

	return PhotoURLOutput{ResourceName: resp.Name, URL: resp.PhotoUri}, nil
}
