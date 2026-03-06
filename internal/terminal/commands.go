package terminal

import (
	"time"

	"github.com/alecthomas/kong"
)

// CLI is the top-level command tree for GoogleMapsCLI.
type CLI struct {
	Globals    GlobalFlags     `embed:""`
	Search     SearchCmd       `cmd:"" help:"Search places by text query."`
	Suggest    SuggestCmd      `cmd:"" help:"Autocomplete places and queries."`
	Nearby     NearbyCmd       `cmd:"" help:"Search nearby places by location."`
	Route      RouteCmd        `cmd:"" help:"Search places along a route."`
	Directions DirectionsCmd   `cmd:"" help:"Get directions and travel time between two points."`
	Info       InfoCmd         `cmd:"" help:"Fetch place details by place ID."`
	Photo      PhotoCmd        `cmd:"" help:"Fetch a photo URL by photo name."`
	Lookup     LookupCmd       `cmd:"" help:"Resolve a location string to candidate places."`
}

// GlobalFlags are shared across all commands.
type GlobalFlags struct {
	APIKey        string        `help:"Google Maps API key." env:"GOOGLE_PLACES_API_KEY"`
	PlacesURL     string        `help:"Places API base URL." env:"GOOGLE_PLACES_BASE_URL" default:"https://places.googleapis.com/v1"`
	RoutesURL     string        `help:"Routes API base URL." env:"GOOGLE_ROUTES_BASE_URL" default:"https://routes.googleapis.com"`
	DirectionsURL string        `help:"Directions API base URL." env:"GOOGLE_DIRECTIONS_BASE_URL" default:"https://routes.googleapis.com"`
	Timeout       time.Duration `help:"HTTP timeout." default:"10s"`
	JSON          bool          `help:"Output as JSON."`
	NoColor       bool          `help:"Disable colored output."`
	Verbose       bool          `help:"Enable verbose logging."`
	Version       VersionFlag   `name:"version" help:"Print version and exit."`
}

type SearchCmd struct {
	Query      string   `arg:"" name:"query" help:"Search text."`
	Limit      int      `help:"Max results (1-20)." default:"10"`
	Cursor     string   `help:"Pagination cursor from previous response."`
	Lang       string   `help:"BCP-47 language code (e.g. en, en-US)."`
	Country    string   `help:"CLDR region code (e.g. US, DE)."`
	Keyword    string   `help:"Keyword to append to the query."`
	Category   []string `help:"Place category filter. Repeatable."`
	OnlyOpen   *bool    `help:"Return only currently open places."`
	MinScore   *float64 `help:"Minimum rating (0-5)."`
	PriceTier  []int    `help:"Price tiers 0-4. Repeatable."`
	Lat        *float64 `help:"Latitude for location bias."`
	Lng        *float64 `help:"Longitude for location bias."`
	Radius     *float64 `help:"Radius in meters for location bias."`
}

type SuggestCmd struct {
	Fragment   string   `arg:"" name:"input" help:"Partial text to complete."`
	Limit      int      `help:"Max suggestions (1-20)." default:"5"`
	Session    string   `help:"Session token for billing consistency."`
	Lang       string   `help:"BCP-47 language code (e.g. en, en-US)."`
	Country    string   `help:"CLDR region code (e.g. US, DE)."`
	Lat        *float64 `help:"Latitude for location bias."`
	Lng        *float64 `help:"Longitude for location bias."`
	Radius     *float64 `help:"Radius in meters for location bias."`
}

type NearbyCmd struct {
	Limit       int      `help:"Max results (1-20)." default:"10"`
	Category    []string `help:"Included place categories. Repeatable."`
	ExcludeCat  []string `help:"Excluded place categories. Repeatable."`
	Lang        string   `help:"BCP-47 language code (e.g. en, en-US)."`
	Country     string   `help:"CLDR region code (e.g. US, DE)."`
	Lat         *float64 `help:"Latitude for location restriction."`
	Lng         *float64 `help:"Longitude for location restriction."`
	Radius      *float64 `help:"Radius in meters for location restriction."`
}

type InfoCmd struct {
	PlaceID  string `arg:"" name:"place_id" help:"Place ID."`
	Lang     string `help:"BCP-47 language code (e.g. en, en-US)."`
	Country  string `help:"CLDR region code (e.g. US, DE)."`
	Reviews  bool   `help:"Include reviews in the response."`
	Photos   bool   `help:"Include photos in the response."`
}

type PhotoCmd struct {
	Name      string `arg:"" name:"photo_name" help:"Photo resource name (places/.../photos/...)."`
	MaxWidth  int    `help:"Max width in pixels." name:"max-width"`
	MaxHeight int    `help:"Max height in pixels." name:"max-height"`
}

type LookupCmd struct {
	Address string `arg:"" name:"location" help:"Location text to resolve."`
	Limit   int    `help:"Max results (1-10)." default:"5"`
	Lang    string `help:"BCP-47 language code (e.g. en, en-US)."`
	Country string `help:"CLDR region code (e.g. US, DE)."`
}

type RouteCmd struct {
	Query    string  `arg:"" name:"query" help:"Search text."`
	From     string  `help:"Origin location (address or place name)."`
	To       string  `help:"Destination location (address or place name)."`
	TravelBy string  `help:"Travel mode: DRIVE, WALK, BICYCLE, TWO_WHEELER, TRANSIT." default:"DRIVE"`
	Radius   float64 `help:"Search radius in meters." default:"1000"`
	Stops    int     `help:"Max sampled stops along the route." default:"5"`
	PerStop  int     `help:"Max results per stop (1-20)." default:"5"`
	Lang     string  `help:"BCP-47 language code (e.g. en, en-US)."`
	Country  string  `help:"CLDR region code (e.g. US, DE)."`
}

type DirectionsCmd struct {
	From        string   `help:"Origin address or place name."`
	To          string   `help:"Destination address or place name."`
	FromID      string   `help:"Origin place ID." name:"from-place-id"`
	ToID        string   `help:"Destination place ID." name:"to-place-id"`
	FromLat     *float64 `help:"Origin latitude." name:"from-lat"`
	FromLng     *float64 `help:"Origin longitude." name:"from-lng"`
	ToLat       *float64 `help:"Destination latitude." name:"to-lat"`
	ToLng       *float64 `help:"Destination longitude." name:"to-lng"`
	TravelBy    string   `help:"Travel mode: walk, drive, bicycle, transit." default:"walk"`
	CompareTo   string   `help:"Compare with another mode: walk, drive, bicycle, transit."`
	ShowSteps   bool     `help:"Include step-by-step instructions." name:"steps"`
	Units       string   `help:"Units: metric or imperial." default:"metric"`
	Lang        string   `help:"BCP-47 language code (e.g. en, en-US)."`
	Country     string   `help:"CLDR region code (e.g. US, DE)."`
}

// VersionFlag prints the version and exits.
type VersionFlag string

func (v VersionFlag) Decode(_ *kong.DecodeContext) error { return nil }
func (v VersionFlag) IsBool() bool                       { return true }
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	_, _ = app.Stdout.Write([]byte(vars["version"] + "\n"))
	app.Exit(0)
	return nil
}

// BuildVersion is set by the linker at build time.
var BuildVersion = "dev"
