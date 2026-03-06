package terminal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/krishnashahane/googlemapscli/pkg/maps"
)

// Handle wires the CLI output and API access.
type Handle struct {
	gm      *maps.GoogleMaps
	out     io.Writer
	errOut  io.Writer
	asJSON  bool
	palette Palette
}

// Execute parses args and runs the appropriate command.
func Execute(args []string, stdout, stderr io.Writer) int {
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}

	root := CLI{}
	exitCode := 0
	parser, err := kong.New(
		&root,
		kong.Name("googlemapscli"),
		kong.Description("GoogleMapsCLI — query Google Maps from your terminal. By Krishna."),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{Compact: true, Summary: true}),
		kong.Writers(stdout, stderr),
		kong.Exit(func(code int) {
			exitCode = code
			panic(haltSignal{code: code})
		}),
		kong.Vars{"version": BuildVersion},
	)
	if err != nil {
		_, _ = fmt.Fprintln(stderr, err)
		return 1
	}

	kctx, halted, err := safeParse(parser, args, &exitCode)
	if halted {
		return exitCode
	}
	if err != nil {
		if pe, ok := err.(*kong.ParseError); ok {
			_ = pe.Context.PrintUsage(true)
			_, _ = fmt.Fprintln(stderr, pe.Error())
			return pe.ExitCode()
		}
		_, _ = fmt.Fprintln(stderr, err)
		return 2
	}
	if root.Globals.JSON {
		root.Globals.NoColor = true
	}

	gm := maps.NewGoogleMaps(maps.Settings{
		Key:           root.Globals.APIKey,
		PlacesURL:     root.Globals.PlacesURL,
		RoutesURL:     root.Globals.RoutesURL,
		DirectionsURL: root.Globals.DirectionsURL,
		Timeout:       root.Globals.Timeout,
	})

	h := &Handle{
		gm:      gm,
		out:     stdout,
		errOut:  stderr,
		asJSON:  root.Globals.JSON,
		palette: NewPalette(ShouldColor(root.Globals.NoColor)),
	}

	kctx.Bind(h)
	if err := kctx.Run(); err != nil {
		return reportError(stderr, err)
	}
	return 0
}

type haltSignal struct{ code int }

func safeParse(parser *kong.Kong, args []string, exitCode *int) (ctx *kong.Context, halted bool, err error) {
	defer func() {
		if r := recover(); r != nil {
			if sig, ok := r.(haltSignal); ok {
				if exitCode != nil {
					*exitCode = sig.code
				}
				halted = true
				ctx = nil
				err = nil
				return
			}
			panic(r)
		}
	}()
	ctx, err = parser.Parse(args)
	return ctx, halted, err
}

// -- Command runners --

func (c *SearchCmd) Run(h *Handle) error {
	in := maps.TextSearchInput{
		Text:       c.Query,
		MaxResults: c.Limit,
		Cursor:     c.Cursor,
		Lang:       c.Lang,
		CountryCode: c.Country,
	}

	f := maps.SearchFilters{}
	hasFilters := false
	if c.Keyword != "" {
		f.Keyword = c.Keyword
		hasFilters = true
	}
	if len(c.Category) > 0 {
		f.Categories = c.Category
		hasFilters = true
	}
	if c.OnlyOpen != nil {
		f.OnlyOpen = c.OnlyOpen
		hasFilters = true
	}
	if c.MinScore != nil {
		f.MinScore = c.MinScore
		hasFilters = true
	}
	if len(c.PriceTier) > 0 {
		f.PriceTiers = c.PriceTier
		hasFilters = true
	}
	if hasFilters {
		in.Filters = &f
	}

	if c.Lat != nil || c.Lng != nil || c.Radius != nil {
		if c.Lat == nil || c.Lng == nil || c.Radius == nil {
			return maps.InputError{Param: "vicinity", Reason: "lat, lng, and radius all required"}
		}
		in.Vicinity = &maps.Area{Latitude: *c.Lat, Longitude: *c.Lng, Radius: *c.Radius}
	}

	result, err := h.gm.TextSearch(context.Background(), in)
	if err != nil {
		return err
	}

	if h.asJSON {
		if err := emitJSON(h.out, result.Places); err != nil {
			return err
		}
		if result.NextCursor != "" {
			_, _ = fmt.Fprintln(h.errOut, "next_cursor:", result.NextCursor)
		}
		return nil
	}

	_, err = fmt.Fprintln(h.out, FormatTextSearch(h.palette, result))
	return err
}

func (c *SuggestCmd) Run(h *Handle) error {
	in := maps.SuggestInput{
		Fragment:    c.Fragment,
		MaxResults:  c.Limit,
		Session:     c.Session,
		Lang:        c.Lang,
		CountryCode: c.Country,
	}

	if c.Lat != nil || c.Lng != nil || c.Radius != nil {
		if c.Lat == nil || c.Lng == nil || c.Radius == nil {
			return maps.InputError{Param: "vicinity", Reason: "lat, lng, and radius all required"}
		}
		in.Vicinity = &maps.Area{Latitude: *c.Lat, Longitude: *c.Lng, Radius: *c.Radius}
	}

	result, err := h.gm.Suggest(context.Background(), in)
	if err != nil {
		return err
	}

	if h.asJSON {
		return emitJSON(h.out, result.Items)
	}

	_, err = fmt.Fprintln(h.out, FormatSuggest(h.palette, result))
	return err
}

func (c *NearbyCmd) Run(h *Handle) error {
	if c.Lat == nil || c.Lng == nil || c.Radius == nil {
		return maps.InputError{Param: "center", Reason: "lat, lng, and radius all required"}
	}

	in := maps.NearbyInput{
		Center:      &maps.Area{Latitude: *c.Lat, Longitude: *c.Lng, Radius: *c.Radius},
		MaxResults:  c.Limit,
		Include:     c.Category,
		Exclude:     c.ExcludeCat,
		Lang:        c.Lang,
		CountryCode: c.Country,
	}

	result, err := h.gm.NearbySearch(context.Background(), in)
	if err != nil {
		return err
	}

	if h.asJSON {
		if err := emitJSON(h.out, result.Places); err != nil {
			return err
		}
		if result.NextCursor != "" {
			_, _ = fmt.Fprintln(h.errOut, "next_cursor:", result.NextCursor)
		}
		return nil
	}

	_, err = fmt.Fprintln(h.out, FormatNearby(h.palette, result))
	return err
}

func (c *InfoCmd) Run(h *Handle) error {
	result, err := h.gm.PlaceDetails(context.Background(), maps.PlaceInfoInput{
		ID:          c.PlaceID,
		Lang:        c.Lang,
		CountryCode: c.Country,
		WithReviews: c.Reviews,
		WithPhotos:  c.Photos,
	})
	if err != nil {
		return err
	}

	if h.asJSON {
		return emitJSON(h.out, result)
	}

	_, err = fmt.Fprintln(h.out, FormatPlaceInfo(h.palette, result))
	return err
}

func (c *PhotoCmd) Run(h *Handle) error {
	result, err := h.gm.FetchPhotoURL(context.Background(), maps.PhotoURLInput{
		ResourceName: c.Name,
		MaxWidth:     c.MaxWidth,
		MaxHeight:    c.MaxHeight,
	})
	if err != nil {
		return err
	}

	if h.asJSON {
		return emitJSON(h.out, result)
	}

	_, err = fmt.Fprintln(h.out, FormatPhotoURL(h.palette, result))
	return err
}

func (c *LookupCmd) Run(h *Handle) error {
	in := maps.LocationLookupInput{
		Address:     c.Address,
		MaxResults:  c.Limit,
		Lang:        c.Lang,
		CountryCode: c.Country,
	}

	result, err := h.gm.LocationLookup(context.Background(), in)
	if err != nil {
		return err
	}

	if h.asJSON {
		return emitJSON(h.out, result.Matches)
	}

	_, err = fmt.Fprintln(h.out, FormatLookup(h.palette, result))
	return err
}

func (c *RouteCmd) Run(h *Handle) error {
	in := maps.RouteSearchInput{
		Text:         c.Query,
		StartPoint:   c.From,
		EndPoint:     c.To,
		TravelBy:     c.TravelBy,
		SearchRadius: c.Radius,
		Stops:        c.Stops,
		PerStop:      c.PerStop,
		Lang:         c.Lang,
		CountryCode:  c.Country,
	}

	result, err := h.gm.RouteSearch(context.Background(), in)
	if err != nil {
		return err
	}

	if h.asJSON {
		return emitJSON(h.out, result)
	}

	_, err = fmt.Fprintln(h.out, FormatRouteSearch(h.palette, result))
	return err
}

func (c *DirectionsCmd) Run(h *Handle) error {
	primary := maps.CanonicalMode(c.TravelBy)
	if primary == "" {
		return maps.InputError{Param: "travel_by", Reason: "must be walk, drive, bicycle, or transit"}
	}
	var secondary string
	if strings.TrimSpace(c.CompareTo) != "" {
		secondary = maps.CanonicalMode(c.CompareTo)
		if secondary == "" {
			return maps.InputError{Param: "compare_to", Reason: "must be walk, drive, bicycle, or transit"}
		}
		if secondary == primary {
			return maps.InputError{Param: "compare_to", Reason: "must differ from travel_by"}
		}
	}

	in := maps.NavigationInput{
		Origin:        c.From,
		Destination:   c.To,
		OriginID:      c.FromID,
		DestinationID: c.ToID,
		TravelBy:      primary,
		MeasureSystem: c.Units,
		Lang:          c.Lang,
		CountryCode:   c.Country,
	}
	if c.FromLat != nil || c.FromLng != nil {
		if c.FromLat == nil || c.FromLng == nil {
			return maps.InputError{Param: "origin_coords", Reason: "both lat and lng required"}
		}
		in.OriginCoords = &maps.Coordinates{Latitude: *c.FromLat, Longitude: *c.FromLng}
	}
	if c.ToLat != nil || c.ToLng != nil {
		if c.ToLat == nil || c.ToLng == nil {
			return maps.InputError{Param: "destination_coords", Reason: "both lat and lng required"}
		}
		in.DestinationCoords = &maps.Coordinates{Latitude: *c.ToLat, Longitude: *c.ToLng}
	}

	result, err := h.gm.Navigate(context.Background(), in)
	if err != nil {
		return err
	}

	var secondResult *maps.NavigationOutput
	if secondary != "" {
		alt := in
		alt.TravelBy = secondary
		altResult, err := h.gm.Navigate(context.Background(), alt)
		if err != nil {
			return err
		}
		secondResult = &altResult
	}

	if h.asJSON {
		if secondResult != nil {
			return emitJSON(h.out, []maps.NavigationOutput{result, *secondResult})
		}
		return emitJSON(h.out, result)
	}

	if secondResult != nil {
		_, err = h.out.Write([]byte(FormatNavigation(h.palette, result, c.ShowSteps)))
		if err != nil {
			return err
		}
		_, err = h.out.Write([]byte("\n\n" + FormatNavigation(h.palette, *secondResult, c.ShowSteps)))
		return err
	}

	_, err = h.out.Write([]byte(FormatNavigation(h.palette, result, c.ShowSteps)))
	return err
}

// -- utilities --

func emitJSON(w io.Writer, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	_, err = w.Write(append(data, '\n'))
	return err
}

func reportError(w io.Writer, err error) int {
	if err == nil {
		return 0
	}
	var inputErr maps.InputError
	if errors.As(err, &inputErr) {
		_, _ = fmt.Fprintln(w, inputErr.Error())
		return 2
	}
	if errors.Is(err, maps.ErrNoAPIKey) {
		_, _ = fmt.Fprintln(w, err.Error())
		return 2
	}
	_, _ = fmt.Fprintln(w, err.Error())
	return 1
}
