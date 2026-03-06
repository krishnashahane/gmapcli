package maps

// wire.go contains internal API payload structs for marshalling/unmarshalling.

type apiSearchResult struct {
	Places        []apiPlace `json:"places"`
	NextPageToken string     `json:"nextPageToken"`
}

type apiPlace struct {
	ID                  string               `json:"id"`
	DisplayName         *apiDisplayName      `json:"displayName,omitempty"`
	FormattedAddress    string               `json:"formattedAddress,omitempty"`
	Location            *apiLocation         `json:"location,omitempty"`
	Rating              *float64             `json:"rating,omitempty"`
	UserRatingCount     *int                 `json:"userRatingCount,omitempty"`
	PriceLevel          string               `json:"priceLevel,omitempty"`
	Types               []string             `json:"types,omitempty"`
	CurrentOpeningHours *apiHours            `json:"currentOpeningHours,omitempty"`
	RegularOpeningHours *apiHours            `json:"regularOpeningHours,omitempty"`
	NationalPhoneNumber string               `json:"nationalPhoneNumber,omitempty"`
	WebsiteUri          string               `json:"websiteUri,omitempty"`
	Reviews             []apiFeedback        `json:"reviews,omitempty"`
	Photos              []apiImage           `json:"photos,omitempty"`
}

type apiDisplayName struct {
	Text string `json:"text"`
}

type apiLocation struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type apiHours struct {
	OpenNow             *bool    `json:"openNow,omitempty"`
	WeekdayDescriptions []string `json:"weekdayDescriptions,omitempty"`
}

type apiFeedback struct {
	Name                           string           `json:"name,omitempty"`
	RelativePublishTimeDescription string           `json:"relativePublishTimeDescription,omitempty"`
	Text                           *apiLangText     `json:"text,omitempty"`
	OriginalText                   *apiLangText     `json:"originalText,omitempty"`
	Rating                         *float64         `json:"rating,omitempty"`
	AuthorAttribution              *apiPerson       `json:"authorAttribution,omitempty"`
	PublishTime                    string           `json:"publishTime,omitempty"`
	FlagContentUri                 string           `json:"flagContentUri,omitempty"`
	GoogleMapsUri                  string           `json:"googleMapsUri,omitempty"`
	VisitDate                      *apiCalDate      `json:"visitDate,omitempty"`
}

type apiLangText struct {
	Text         string `json:"text,omitempty"`
	LanguageCode string `json:"languageCode,omitempty"`
}

type apiPerson struct {
	DisplayName string `json:"displayName,omitempty"`
	URI         string `json:"uri,omitempty"`
	PhotoURI    string `json:"photoUri,omitempty"`
}

type apiCalDate struct {
	Year  int `json:"year,omitempty"`
	Month int `json:"month,omitempty"`
	Day   int `json:"day,omitempty"`
}

type apiImage struct {
	Name               string      `json:"name,omitempty"`
	WidthPx            int         `json:"widthPx,omitempty"`
	HeightPx           int         `json:"heightPx,omitempty"`
	AuthorAttributions []apiPerson `json:"authorAttributions,omitempty"`
}

type apiPhotoMedia struct {
	Name     string `json:"name,omitempty"`
	PhotoUri string `json:"photoUri,omitempty"`
}

type apiSuggestResult struct {
	Suggestions []apiSuggestionItem `json:"suggestions"`
}

type apiSuggestionItem struct {
	PlacePrediction *apiPlacePrediction `json:"placePrediction,omitempty"`
	QueryPrediction *apiQueryPrediction `json:"queryPrediction,omitempty"`
}

type apiPlacePrediction struct {
	PlaceId          string             `json:"placeId,omitempty"`
	Place            string             `json:"place,omitempty"`
	Text             *apiSmallText      `json:"text,omitempty"`
	StructuredFormat *apiStructuredText `json:"structuredFormat,omitempty"`
	Types            []string           `json:"types,omitempty"`
	DistanceMeters   *int               `json:"distanceMeters,omitempty"`
}

type apiQueryPrediction struct {
	Text             *apiSmallText      `json:"text,omitempty"`
	StructuredFormat *apiStructuredText `json:"structuredFormat,omitempty"`
}

type apiStructuredText struct {
	MainText      *apiSmallText `json:"mainText,omitempty"`
	SecondaryText *apiSmallText `json:"secondaryText,omitempty"`
}

type apiSmallText struct {
	Text string `json:"text,omitempty"`
}

type apiRoutesResult struct {
	Routes []apiRouteEntry `json:"routes"`
}

type apiRouteEntry struct {
	Polyline    apiPolyline     `json:"polyline"`
	Description string          `json:"description,omitempty"`
	Warnings    []string        `json:"warnings,omitempty"`
	Legs        []apiRouteLeg   `json:"legs"`
}

type apiPolyline struct {
	EncodedPolyline string `json:"encodedPolyline"`
}

type apiRouteLeg struct {
	DistanceMeters  int                `json:"distanceMeters"`
	Duration        string             `json:"duration"`
	LocalizedValues apiLegLocalized    `json:"localizedValues"`
	Steps           []apiRouteStep     `json:"steps"`
}

type apiRouteStep struct {
	DistanceMeters        int                   `json:"distanceMeters"`
	StaticDuration        string                `json:"staticDuration"`
	TravelMode            string                `json:"travelMode,omitempty"`
	NavigationInstruction apiNavInstruction     `json:"navigationInstruction"`
	LocalizedValues       apiStepLocalized      `json:"localizedValues"`
}

type apiNavInstruction struct {
	Instructions string `json:"instructions,omitempty"`
	Maneuver     string `json:"maneuver,omitempty"`
}

type apiLegLocalized struct {
	Distance apiLabelText `json:"distance"`
	Duration apiLabelText `json:"duration"`
}

type apiStepLocalized struct {
	Distance       apiLabelText `json:"distance"`
	StaticDuration apiLabelText `json:"staticDuration"`
}

type apiLabelText struct {
	Text string `json:"text,omitempty"`
}
