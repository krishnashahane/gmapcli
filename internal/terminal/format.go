package terminal

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/krishnashahane/googlemapscli/pkg/maps"
)

const noData = "No results."

func FormatTextSearch(p Palette, out maps.TextSearchOutput) string {
	var buf bytes.Buffer
	n := len(out.Places)
	if n == 0 {
		return noData
	}
	buf.WriteString(p.Strong(fmt.Sprintf("Results (%d)", n)))
	buf.WriteByte('\n')

	for i, place := range out.Places {
		buf.WriteString(fmt.Sprintf("%d. %s\n", i+1, heading(p, place.Title, place.FullAddress)))
		printBrief(&buf, p, place)
		if i < n-1 {
			buf.WriteByte('\n')
		}
	}

	if strings.TrimSpace(out.NextCursor) != "" {
		buf.WriteByte('\n')
		buf.WriteString(p.Muted("Next cursor:"))
		buf.WriteByte(' ')
		buf.WriteString(out.NextCursor)
	}

	return buf.String()
}

func FormatSuggest(p Palette, out maps.SuggestOutput) string {
	var buf bytes.Buffer
	n := len(out.Items)
	if n == 0 {
		return noData
	}
	buf.WriteString(p.Strong(fmt.Sprintf("Suggestions (%d)", n)))
	buf.WriteByte('\n')

	for i, item := range out.Items {
		title := item.Primary
		if strings.TrimSpace(title) == "" {
			title = item.Label
		}
		sub := item.Secondary
		if strings.TrimSpace(sub) == "" && strings.TrimSpace(item.Label) != "" && strings.TrimSpace(item.Primary) != "" {
			sub = item.Label
		}
		buf.WriteString(fmt.Sprintf("%d. %s\n", i+1, heading(p, title, sub)))
		field(&buf, p, "Type", item.Type)
		field(&buf, p, "ID", item.PlaceID)
		field(&buf, p, "Resource", item.Resource)
		printCategories(&buf, p, item.Categories)
		if item.DistanceM != nil {
			field(&buf, p, "Distance", fmt.Sprintf("%dm", *item.DistanceM))
		}
		if i < n-1 {
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}

func FormatNearby(p Palette, out maps.NearbyOutput) string {
	var buf bytes.Buffer
	n := len(out.Places)
	if n == 0 {
		return noData
	}
	buf.WriteString(p.Strong(fmt.Sprintf("Nearby (%d)", n)))
	buf.WriteByte('\n')

	for i, place := range out.Places {
		buf.WriteString(fmt.Sprintf("%d. %s\n", i+1, heading(p, place.Title, place.FullAddress)))
		printBrief(&buf, p, place)
		if i < n-1 {
			buf.WriteByte('\n')
		}
	}

	if strings.TrimSpace(out.NextCursor) != "" {
		buf.WriteByte('\n')
		buf.WriteString(p.Muted("Next cursor:"))
		buf.WriteByte(' ')
		buf.WriteString(out.NextCursor)
	}

	return buf.String()
}

func FormatPlaceInfo(p Palette, info maps.PlaceInfo) string {
	var buf bytes.Buffer
	buf.WriteString(p.Strong(heading(p, info.Title, info.FullAddress)))
	buf.WriteByte('\n')
	field(&buf, p, "ID", info.ID)
	printCoords(&buf, p, info.Coords)
	printScoring(&buf, p, info.Score, info.Votes, info.PriceTier)
	printCategories(&buf, p, info.Categories)
	printOpenStatus(&buf, p, info.IsOpen)
	field(&buf, p, "Phone", info.PhoneNumber)
	field(&buf, p, "Website", info.Homepage)
	printImages(&buf, p, info.Images)
	printFeedback(&buf, p, info.Feedback)
	if len(info.Schedule) > 0 {
		buf.WriteString(p.Muted("Hours:"))
		buf.WriteByte('\n')
		for _, line := range info.Schedule {
			buf.WriteString("  - ")
			buf.WriteString(line)
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}

func FormatPhotoURL(p Palette, out maps.PhotoURLOutput) string {
	var buf bytes.Buffer
	buf.WriteString(p.Strong("Photo"))
	buf.WriteByte('\n')
	field(&buf, p, "Name", out.ResourceName)
	field(&buf, p, "URL", out.URL)
	return buf.String()
}

func FormatLookup(p Palette, out maps.LocationLookupOutput) string {
	var buf bytes.Buffer
	n := len(out.Matches)
	if n == 0 {
		return noData
	}
	buf.WriteString(p.Strong(fmt.Sprintf("Resolved (%d)", n)))
	buf.WriteByte('\n')

	for i, m := range out.Matches {
		buf.WriteString(fmt.Sprintf("%d. %s\n", i+1, heading(p, m.Title, m.FullAddress)))
		field(&buf, p, "ID", m.ID)
		printCoords(&buf, p, m.Coords)
		printCategories(&buf, p, m.Categories)
		if i < n-1 {
			buf.WriteByte('\n')
		}
	}
	return buf.String()
}

func FormatRouteSearch(p Palette, out maps.RouteSearchOutput) string {
	var buf bytes.Buffer
	n := len(out.Stops)
	if n == 0 {
		return noData
	}
	buf.WriteString(p.Strong(fmt.Sprintf("Route stops (%d)", n)))
	buf.WriteByte('\n')

	for i, stop := range out.Stops {
		buf.WriteString(p.Strong(fmt.Sprintf("Stop %d", i+1)))
		buf.WriteByte(' ')
		buf.WriteString(p.Muted(fmt.Sprintf("(%.6f, %.6f)", stop.Point.Latitude, stop.Point.Longitude)))
		buf.WriteByte('\n')

		if len(stop.Places) == 0 {
			buf.WriteString(noData)
			buf.WriteByte('\n')
		} else {
			for j, place := range stop.Places {
				buf.WriteString(fmt.Sprintf("%d. %s\n", j+1, heading(p, place.Title, place.FullAddress)))
				printBrief(&buf, p, place)
				if j < len(stop.Places)-1 {
					buf.WriteByte('\n')
				}
			}
		}

		if i < n-1 {
			buf.WriteByte('\n')
		}
	}

	return buf.String()
}

func FormatNavigation(p Palette, out maps.NavigationOutput, showSteps bool) string {
	var buf bytes.Buffer
	label := "Directions"
	if m := strings.TrimSpace(out.TravelBy); m != "" {
		label = fmt.Sprintf("Directions (%s)", m)
	}
	buf.WriteString(p.Strong(label))
	buf.WriteByte('\n')
	field(&buf, p, "From", out.FromLabel)
	field(&buf, p, "To", out.ToLabel)
	field(&buf, p, "Summary", out.Description)
	field(&buf, p, "Distance", out.DistanceLabel)
	field(&buf, p, "Duration", out.DurationLabel)
	if len(out.Alerts) > 0 {
		buf.WriteString(p.Muted("Warnings:"))
		buf.WriteByte('\n')
		for _, a := range out.Alerts {
			if strings.TrimSpace(a) == "" {
				continue
			}
			buf.WriteString("  - ")
			buf.WriteString(a)
			buf.WriteByte('\n')
		}
	}
	if showSteps {
		buf.WriteString(p.Muted("Steps:"))
		buf.WriteByte('\n')
		if len(out.Maneuvers) == 0 {
			buf.WriteString("  - ")
			buf.WriteString(noData)
			buf.WriteByte('\n')
		} else {
			for i, step := range out.Maneuvers {
				line := stepText(step)
				if line == "" {
					continue
				}
				buf.WriteString(fmt.Sprintf("  %d. %s\n", i+1, line))
			}
		}
	}
	return buf.String()
}

// -- helpers --

func heading(p Palette, name, address string) string {
	display := strings.TrimSpace(name)
	if display == "" {
		display = "(unnamed)"
	}
	if address == "" {
		return p.Teal(display)
	}
	return p.Teal(display) + " — " + address
}

func field(buf *bytes.Buffer, p Palette, label, value string) {
	if strings.TrimSpace(value) == "" {
		return
	}
	buf.WriteString(p.Muted(label + ":"))
	buf.WriteByte(' ')
	buf.WriteString(value)
	buf.WriteByte('\n')
}

func printBrief(buf *bytes.Buffer, p Palette, place maps.PlaceBrief) {
	field(buf, p, "ID", place.ID)
	printCoords(buf, p, place.Coords)
	printScoring(buf, p, place.Score, place.Votes, place.PriceTier)
	printCategories(buf, p, place.Categories)
	printOpenStatus(buf, p, place.IsOpen)
}

func printCoords(buf *bytes.Buffer, p Palette, c *maps.Coordinates) {
	if c == nil {
		return
	}
	field(buf, p, "Location", fmt.Sprintf("%.6f, %.6f", c.Latitude, c.Longitude))
}

func printScoring(buf *bytes.Buffer, p Palette, score *float64, votes *int, tier *int) {
	if score == nil && votes == nil && tier == nil {
		return
	}
	parts := make([]string, 0, 3)
	if score != nil {
		s := fmt.Sprintf("%.1f", *score)
		if votes != nil {
			s += fmt.Sprintf(" (%d)", *votes)
		}
		parts = append(parts, s)
	} else if votes != nil {
		parts = append(parts, fmt.Sprintf("%d ratings", *votes))
	}
	if tier != nil {
		parts = append(parts, fmt.Sprintf("$%d", *tier))
	}
	field(buf, p, "Rating", strings.Join(parts, " · "))
}

func printCategories(buf *bytes.Buffer, p Palette, cats []string) {
	if len(cats) == 0 {
		return
	}
	unique := dedupSort(cats)
	field(buf, p, "Types", strings.Join(unique, ", "))
}

func printOpenStatus(buf *bytes.Buffer, p Palette, open *bool) {
	if open == nil {
		return
	}
	v := "no"
	if *open {
		v = "yes"
	}
	field(buf, p, "Open now", v)
}

func printImages(buf *bytes.Buffer, p Palette, images []maps.ImageRef) {
	if len(images) == 0 {
		return
	}
	buf.WriteString(p.Muted("Photos:"))
	buf.WriteByte('\n')
	cap := 3
	if len(images) < cap {
		cap = len(images)
	}
	for i := 0; i < cap; i++ {
		img := images[i]
		parts := make([]string, 0, 3)
		if strings.TrimSpace(img.ResourceName) != "" {
			parts = append(parts, img.ResourceName)
		}
		if img.Width > 0 && img.Height > 0 {
			parts = append(parts, fmt.Sprintf("%dx%d", img.Width, img.Height))
		}
		if len(img.Credits) > 0 && strings.TrimSpace(img.Credits[0].Name) != "" {
			parts = append(parts, "by "+img.Credits[0].Name)
		}
		if len(parts) > 0 {
			buf.WriteString("  - ")
			buf.WriteString(strings.Join(parts, " · "))
			buf.WriteByte('\n')
		}
	}
	if len(images) > 3 {
		buf.WriteString(p.Muted(fmt.Sprintf("  ... %d more", len(images)-3)))
		buf.WriteByte('\n')
	}
}

func printFeedback(buf *bytes.Buffer, p Palette, items []maps.FeedbackEntry) {
	if len(items) == 0 {
		return
	}
	buf.WriteString(p.Muted("Reviews:"))
	buf.WriteByte('\n')
	cap := 3
	if len(items) < cap {
		cap = len(items)
	}
	for i := 0; i < cap; i++ {
		fb := items[i]
		parts := make([]string, 0, 4)
		if fb.Stars != nil {
			parts = append(parts, fmt.Sprintf("%.1f stars", *fb.Stars))
		}
		if fb.Reviewer != nil && strings.TrimSpace(fb.Reviewer.Name) != "" {
			parts = append(parts, "by "+fb.Reviewer.Name)
		}
		if strings.TrimSpace(fb.TimeAgo) != "" {
			parts = append(parts, "("+fb.TimeAgo+")")
		}
		text := feedbackBody(fb)
		if text != "" {
			parts = append(parts, text)
		}
		if len(parts) > 0 {
			buf.WriteString("  - ")
			buf.WriteString(strings.Join(parts, " "))
			buf.WriteByte('\n')
		}
	}
	if len(items) > 3 {
		buf.WriteString(p.Muted(fmt.Sprintf("  ... %d more", len(items)-3)))
		buf.WriteByte('\n')
	}
}

func feedbackBody(fb maps.FeedbackEntry) string {
	text := ""
	if fb.Body != nil {
		text = fb.Body.Content
	}
	if strings.TrimSpace(text) == "" && fb.OrigBody != nil {
		text = fb.OrigBody.Content
	}
	text = strings.TrimSpace(text)
	if len(text) > 200 {
		text = strings.TrimSpace(text[:200]) + "..."
	}
	return text
}

func stepText(s maps.NavStep) string {
	dir := strings.TrimSpace(s.Direction)
	if dir == "" {
		dir = "(no instruction)"
	}
	parts := []string{dir}
	if strings.TrimSpace(s.DistanceLabel) != "" {
		parts = append(parts, s.DistanceLabel)
	}
	if strings.TrimSpace(s.DurationLabel) != "" {
		parts = append(parts, s.DurationLabel)
	}
	return strings.Join(parts, " · ")
}

func dedupSort(vals []string) []string {
	seen := make(map[string]struct{}, len(vals))
	out := make([]string, 0, len(vals))
	for _, v := range vals {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if _, exists := seen[v]; exists {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	sort.Strings(out)
	return out
}
