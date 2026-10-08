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
	name := strings.TrimSpace(in.ResourceName)
	if name == "" || len(name) > 1024 || strings.ContainsAny(name, "\r\n") || !strings.HasPrefix(name, "places/") || !strings.Contains(name, "/photos/") {
		return PhotoURLOutput{}, InputError{Param: "resource_name", Reason: "must be a valid places/.../photos/... resource name"}
	}

	path := "/" + strings.TrimPrefix(name, "/") + "/media"
	params := map[string]string{"skipHttpRedirect": "true"}
	if in.MaxWidth < 0 || in.MaxWidth > 4800 || in.MaxHeight < 0 || in.MaxHeight > 4800 {
		return PhotoURLOutput{}, InputError{Param: "size", Reason: "max width/height must be between 0 and 4800"}
	}
	if in.MaxWidth > 0 { params["maxWidthPx"] = strconv.Itoa(in.MaxWidth) }
	if in.MaxHeight > 0 { params["maxHeightPx"] = strconv.Itoa(in.MaxHeight) }

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
