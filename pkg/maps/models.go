package maps

// Coordinates holds a geographic point.
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Area defines a circular geographic region.
type Area struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    float64 `json:"radius"`
}

// SearchFilters refine text search queries.
type SearchFilters struct {
	Keyword     string   `json:"keyword,omitempty"`
	Categories  []string `json:"categories,omitempty"`
	OnlyOpen    *bool    `json:"only_open,omitempty"`
	MinScore    *float64 `json:"min_score,omitempty"`
	PriceTiers  []int    `json:"price_tiers,omitempty"`
}

// -- Requests --

type TextSearchInput struct {
	Text        string         `json:"text"`
	Filters     *SearchFilters `json:"filters,omitempty"`
	Vicinity    *Area          `json:"vicinity,omitempty"`
	MaxResults  int            `json:"max_results,omitempty"`
	Cursor      string         `json:"cursor,omitempty"`
	Lang        string         `json:"lang,omitempty"`
	CountryCode string         `json:"country_code,omitempty"`
}

type SuggestInput struct {
	Fragment    string `json:"fragment"`
	Session     string `json:"session,omitempty"`
	MaxResults  int    `json:"max_results,omitempty"`
	Lang        string `json:"lang,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
	Vicinity    *Area  `json:"vicinity,omitempty"`
}

type NearbyInput struct {
	Center      *Area    `json:"center,omitempty"`
	MaxResults  int      `json:"max_results,omitempty"`
	Include     []string `json:"include,omitempty"`
	Exclude     []string `json:"exclude,omitempty"`
	Lang        string   `json:"lang,omitempty"`
	CountryCode string   `json:"country_code,omitempty"`
}

type PlaceInfoInput struct {
	ID          string `json:"id"`
	Lang        string `json:"lang,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
	WithReviews bool   `json:"with_reviews,omitempty"`
	WithPhotos  bool   `json:"with_photos,omitempty"`
}

type PhotoURLInput struct {
	ResourceName string `json:"resource_name"`
	MaxWidth     int    `json:"max_width,omitempty"`
	MaxHeight    int    `json:"max_height,omitempty"`
}

type LocationLookupInput struct {
	Address    string `json:"address"`
	MaxResults int    `json:"max_results,omitempty"`
	Lang       string `json:"lang,omitempty"`
	CountryCode string `json:"country_code,omitempty"`
}

type NavigationInput struct {
	Origin          string       `json:"origin,omitempty"`
	Destination     string       `json:"destination,omitempty"`
	OriginID        string       `json:"origin_id,omitempty"`
	DestinationID   string       `json:"destination_id,omitempty"`
	OriginCoords    *Coordinates `json:"origin_coords,omitempty"`
	DestinationCoords *Coordinates `json:"destination_coords,omitempty"`
	TravelBy        string       `json:"travel_by,omitempty"`
	Lang            string       `json:"lang,omitempty"`
	CountryCode     string       `json:"country_code,omitempty"`
	MeasureSystem   string       `json:"measure_system,omitempty"`
}

type RouteSearchInput struct {
	Text         string  `json:"text"`
	StartPoint   string  `json:"start_point"`
	EndPoint     string  `json:"end_point"`
	TravelBy     string  `json:"travel_by,omitempty"`
	SearchRadius float64 `json:"search_radius,omitempty"`
	Stops        int     `json:"stops,omitempty"`
	PerStop      int     `json:"per_stop,omitempty"`
	Lang         string  `json:"lang,omitempty"`
	CountryCode  string  `json:"country_code,omitempty"`
}

// -- Responses --

type TextSearchOutput struct {
	Places     []PlaceBrief `json:"places"`
	NextCursor string       `json:"next_cursor,omitempty"`
}

type NearbyOutput struct {
	Places     []PlaceBrief `json:"places"`
	NextCursor string       `json:"next_cursor,omitempty"`
}

type SuggestOutput struct {
	Items []Suggestion `json:"items"`
}

type PlaceBrief struct {
	ID          string       `json:"id"`
	Title       string       `json:"title,omitempty"`
	FullAddress string       `json:"full_address,omitempty"`
	Coords      *Coordinates `json:"coords,omitempty"`
	Score       *float64     `json:"score,omitempty"`
	Votes       *int         `json:"votes,omitempty"`
	PriceTier   *int         `json:"price_tier,omitempty"`
	Categories  []string     `json:"categories,omitempty"`
	IsOpen      *bool        `json:"is_open,omitempty"`
}

type PlaceInfo struct {
	ID          string       `json:"id"`
	Title       string       `json:"title,omitempty"`
	FullAddress string       `json:"full_address,omitempty"`
	Coords      *Coordinates `json:"coords,omitempty"`
	Score       *float64     `json:"score,omitempty"`
	Votes       *int         `json:"votes,omitempty"`
	PriceTier   *int         `json:"price_tier,omitempty"`
	Categories  []string     `json:"categories,omitempty"`
	IsOpen      *bool        `json:"is_open,omitempty"`
	PhoneNumber string       `json:"phone_number,omitempty"`
	Homepage    string       `json:"homepage,omitempty"`
	Schedule    []string     `json:"schedule,omitempty"`
	Feedback    []FeedbackEntry `json:"feedback,omitempty"`
	Images      []ImageRef   `json:"images,omitempty"`
}

type Suggestion struct {
	Type       string   `json:"type"`
	PlaceID    string   `json:"place_id,omitempty"`
	Resource   string   `json:"resource,omitempty"`
	Label      string   `json:"label,omitempty"`
	Primary    string   `json:"primary,omitempty"`
	Secondary  string   `json:"secondary,omitempty"`
	Categories []string `json:"categories,omitempty"`
	DistanceM  *int     `json:"distance_m,omitempty"`
}

type FeedbackEntry struct {
	RefName     string          `json:"ref_name,omitempty"`
	TimeAgo     string          `json:"time_ago,omitempty"`
	Body        *TranslatedText `json:"body,omitempty"`
	OrigBody    *TranslatedText `json:"orig_body,omitempty"`
	Stars       *float64        `json:"stars,omitempty"`
	Reviewer    *PersonRef      `json:"reviewer,omitempty"`
	PublishedAt string          `json:"published_at,omitempty"`
	FlagURL     string          `json:"flag_url,omitempty"`
	MapLink     string          `json:"map_link,omitempty"`
	VisitedOn   *CalendarDate   `json:"visited_on,omitempty"`
}

type TranslatedText struct {
	Content  string `json:"content,omitempty"`
	LangCode string `json:"lang_code,omitempty"`
}

type PersonRef struct {
	Name     string `json:"name,omitempty"`
	Profile  string `json:"profile,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
}

type CalendarDate struct {
	Year  int `json:"year,omitempty"`
	Month int `json:"month,omitempty"`
	Day   int `json:"day,omitempty"`
}

type ImageRef struct {
	ResourceName string      `json:"resource_name,omitempty"`
	Width        int         `json:"width,omitempty"`
	Height       int         `json:"height,omitempty"`
	Credits      []PersonRef `json:"credits,omitempty"`
}

type PhotoURLOutput struct {
	ResourceName string `json:"resource_name,omitempty"`
	URL          string `json:"url,omitempty"`
}

type LocationMatch struct {
	ID          string       `json:"id"`
	Title       string       `json:"title,omitempty"`
	FullAddress string       `json:"full_address,omitempty"`
	Coords      *Coordinates `json:"coords,omitempty"`
	Categories  []string     `json:"categories,omitempty"`
}

type LocationLookupOutput struct {
	Matches []LocationMatch `json:"matches"`
}

type NavigationOutput struct {
	TravelBy       string          `json:"travel_by"`
	Description    string          `json:"description,omitempty"`
	FromLabel      string          `json:"from_label,omitempty"`
	ToLabel        string          `json:"to_label,omitempty"`
	DistanceLabel  string          `json:"distance_label,omitempty"`
	DistanceM      int             `json:"distance_m,omitempty"`
	DurationLabel  string          `json:"duration_label,omitempty"`
	DurationSec    int             `json:"duration_sec,omitempty"`
	Alerts         []string        `json:"alerts,omitempty"`
	Maneuvers      []NavStep       `json:"maneuvers,omitempty"`
}

type NavStep struct {
	Direction     string `json:"direction,omitempty"`
	DistanceLabel string `json:"distance_label,omitempty"`
	DistanceM     int    `json:"distance_m,omitempty"`
	DurationLabel string `json:"duration_label,omitempty"`
	DurationSec   int    `json:"duration_sec,omitempty"`
	Method        string `json:"method,omitempty"`
	Action        string `json:"action,omitempty"`
}

type RouteSearchOutput struct {
	Stops []RouteStop `json:"stops"`
}

type RouteStop struct {
	Point  Coordinates  `json:"point"`
	Places []PlaceBrief `json:"places"`
}
