package maps

const (
	PlacesEndpoint     = "https://places.googleapis.com/v1"
	RoutesEndpoint     = "https://routes.googleapis.com"
	DirectionsEndpoint = "https://routes.googleapis.com"
)

const (
	SearchMaxResults       = 20
	SearchDefaultResults   = 10
	NearbyMaxResults       = 20
	NearbyDefaultResults   = 10
	ResolveMaxResults      = 10
	ResolveDefaultResults  = 5
	SuggestMaxResults      = 20
	SuggestDefaultResults  = 5
)

const (
	RouteDefaultSamples    = 5
	RouteMaxSamples        = 20
	RouteDefaultSearchArea = 1000
	RouteDefaultPerStop    = 5
)

const (
	PriceFree          = "PRICE_LEVEL_FREE"
	PriceCheap         = "PRICE_LEVEL_INEXPENSIVE"
	PriceMid           = "PRICE_LEVEL_MODERATE"
	PriceHigh          = "PRICE_LEVEL_EXPENSIVE"
	PricePremium       = "PRICE_LEVEL_VERY_EXPENSIVE"
)

var numericToPrice = map[int]string{
	0: PriceFree,
	1: PriceCheap,
	2: PriceMid,
	3: PriceHigh,
	4: PricePremium,
}

var priceToNumeric = map[string]int{
	PriceFree:    0,
	PriceCheap:   1,
	PriceMid:     2,
	PriceHigh:    3,
	PricePremium: 4,
}
