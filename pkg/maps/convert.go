package maps

import "strings"

func toBrief(p apiPlace) PlaceBrief {
	return PlaceBrief{
		ID:          p.ID,
		Title:       extractTitle(p.DisplayName),
		FullAddress: p.FormattedAddress,
		Coords:      toCoords(p.Location),
		Score:       p.Rating,
		Votes:       p.UserRatingCount,
		PriceTier:   toPriceTier(p.PriceLevel),
		Categories:  p.Types,
		IsOpen:      extractOpen(p.CurrentOpeningHours),
	}
}

func toPlaceInfo(p apiPlace) PlaceInfo {
	return PlaceInfo{
		ID:          p.ID,
		Title:       extractTitle(p.DisplayName),
		FullAddress: p.FormattedAddress,
		Coords:      toCoords(p.Location),
		Score:       p.Rating,
		Votes:       p.UserRatingCount,
		PriceTier:   toPriceTier(p.PriceLevel),
		Categories:  p.Types,
		PhoneNumber: p.NationalPhoneNumber,
		Homepage:    p.WebsiteUri,
		Schedule:    extractSchedule(p.RegularOpeningHours),
		IsOpen:      extractOpen(p.CurrentOpeningHours),
		Feedback:    toFeedbackList(p.Reviews),
		Images:      toImageList(p.Photos),
	}
}

func toLocationMatch(p apiPlace) LocationMatch {
	return LocationMatch{
		ID:          p.ID,
		Title:       extractTitle(p.DisplayName),
		FullAddress: p.FormattedAddress,
		Coords:      toCoords(p.Location),
		Categories:  p.Types,
	}
}

func toCoords(loc *apiLocation) *Coordinates {
	if loc == nil {
		return nil
	}
	return &Coordinates{Latitude: loc.Latitude, Longitude: loc.Longitude}
}

func extractTitle(dn *apiDisplayName) string {
	if dn == nil {
		return ""
	}
	return dn.Text
}

func extractOpen(hrs *apiHours) *bool {
	if hrs == nil {
		return nil
	}
	return hrs.OpenNow
}

func extractSchedule(hrs *apiHours) []string {
	if hrs == nil {
		return nil
	}
	return hrs.WeekdayDescriptions
}

func toPriceTier(raw string) *int {
	if raw == "" {
		return nil
	}
	if val, found := priceToNumeric[raw]; found {
		return &val
	}
	return nil
}

func toFeedbackList(items []apiFeedback) []FeedbackEntry {
	if len(items) == 0 {
		return nil
	}
	out := make([]FeedbackEntry, 0, len(items))
	for _, item := range items {
		out = append(out, FeedbackEntry{
			RefName:     item.Name,
			TimeAgo:     item.RelativePublishTimeDescription,
			Body:        toTranslated(item.Text),
			OrigBody:    toTranslated(item.OriginalText),
			Stars:       item.Rating,
			Reviewer:    toPersonRef(item.AuthorAttribution),
			PublishedAt: item.PublishTime,
			FlagURL:     item.FlagContentUri,
			MapLink:     item.GoogleMapsUri,
			VisitedOn:   toCalDate(item.VisitDate),
		})
	}
	return out
}

func toTranslated(t *apiLangText) *TranslatedText {
	if t == nil {
		return nil
	}
	if strings.TrimSpace(t.Text) == "" && strings.TrimSpace(t.LanguageCode) == "" {
		return nil
	}
	return &TranslatedText{Content: t.Text, LangCode: t.LanguageCode}
}

func toPersonRef(p *apiPerson) *PersonRef {
	if p == nil {
		return nil
	}
	if strings.TrimSpace(p.DisplayName) == "" && strings.TrimSpace(p.URI) == "" && strings.TrimSpace(p.PhotoURI) == "" {
		return nil
	}
	return &PersonRef{Name: p.DisplayName, Profile: p.URI, Avatar: p.PhotoURI}
}

func toPersonRefs(items []apiPerson) []PersonRef {
	if len(items) == 0 {
		return nil
	}
	out := make([]PersonRef, 0, len(items))
	for _, item := range items {
		out = append(out, PersonRef{Name: item.DisplayName, Profile: item.URI, Avatar: item.PhotoURI})
	}
	return out
}

func toCalDate(d *apiCalDate) *CalendarDate {
	if d == nil {
		return nil
	}
	if d.Year == 0 && d.Month == 0 && d.Day == 0 {
		return nil
	}
	return &CalendarDate{Year: d.Year, Month: d.Month, Day: d.Day}
}

func toImageList(items []apiImage) []ImageRef {
	if len(items) == 0 {
		return nil
	}
	out := make([]ImageRef, 0, len(items))
	for _, item := range items {
		out = append(out, ImageRef{
			ResourceName: item.Name,
			Width:        item.WidthPx,
			Height:       item.HeightPx,
			Credits:      toPersonRefs(item.AuthorAttributions),
		})
	}
	return out
}

func areaToCircle(a *Area) map[string]any {
	return map[string]any{
		"circle": map[string]any{
			"center": map[string]any{
				"latitude":  a.Latitude,
				"longitude": a.Longitude,
			},
			"radius": a.Radius,
		},
	}
}
