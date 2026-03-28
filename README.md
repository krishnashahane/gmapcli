# GoogleMapsCLI

A fast, modern CLI and Go library for the Google Maps Platform (Places API New + Routes API). Built by Krishna Shahane.

## Features

- Text search with filters: keyword, category, open now, min rating, price tiers
- Autocomplete suggestions for places and queries (session token support)
- Nearby search around a geographic point
- Place details: hours, phone, website, rating, reviews, photos
- Photo media URLs from photo resource names
- Route search along a driving/walking path (Routes API)
- Directions between two points with distance, duration, and turn-by-turn steps
- Location lookup: resolve freeform addresses to place candidates
- Location bias and pagination support
- Color terminal output + `--json` mode (respects `NO_COLOR`)

## Install

```bash
go install github.com/krishnashahane/googlemapscli/cmd/googlemapscli@latest
```

Or build from source:

```bash
make build
```

## Setup

```bash
export GOOGLE_PLACES_API_KEY="your-key-here"
```

Enable the **Places API (New)** and **Routes API** in your Google Cloud Console.

## Usage

```
googlemapscli [flags] <command>

Commands:
  search      Search places by text query
  suggest     Autocomplete places and queries
  nearby      Search nearby places by location
  route       Search places along a route
  directions  Get directions between two points
  info        Fetch place details by place ID
  photo       Fetch a photo URL by photo name
  lookup      Resolve a location string to candidate places
```

### Examples

```bash
googlemapscli search "coffee" --min-score 4 --only-open --limit 5 \
  --lat 40.8065 --lng -73.9719 --radius 3000

googlemapscli suggest "cof" --session "my-session" --limit 5

googlemapscli nearby --lat 47.6062 --lng -122.3321 --radius 1500 --category cafe

googlemapscli route "coffee" --from "Seattle, WA" --to "Portland, OR" --stops 5

googlemapscli directions --from "Pike Place Market" --to "Space Needle" --steps

googlemapscli info ChIJN1t_tDeuEmsRUsoyG83frY4 --reviews --photos

googlemapscli lookup "Riverside Park, New York" --limit 5

googlemapscli search "sushi" --json
```

## Library Usage

```go
gm := maps.NewGoogleMaps(maps.Settings{
    Key:     os.Getenv("GOOGLE_PLACES_API_KEY"),
    Timeout: 8 * time.Second,
})

result, err := gm.TextSearch(ctx, maps.TextSearchInput{
    Text:       "italian restaurant",
    MaxResults: 10,
    Vicinity:   &maps.Area{Latitude: 40.8065, Longitude: -73.9719, Radius: 3000},
})

info, err := gm.PlaceDetails(ctx, maps.PlaceInfoInput{
    ID:          "ChIJN1t_tDeuEmsRUsoyG83frY4",
    WithReviews: true,
})

nav, err := gm.Navigate(ctx, maps.NavigationInput{
    Origin:      "Pike Place Market",
    Destination: "Space Needle",
    TravelBy:    "walking",
})
```

## Testing

```bash
make test
make lint
make coverage
```

## License

MIT License
