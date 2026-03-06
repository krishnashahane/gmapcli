package maps

import "fmt"

func checkArea(a *Area) error {
	if a == nil {
		return nil
	}
	if a.Radius <= 0 {
		return InputError{Param: "area.radius", Reason: "must be positive"}
	}
	if a.Latitude < -90 || a.Latitude > 90 {
		return InputError{Param: "area.latitude", Reason: "must be between -90 and 90"}
	}
	if a.Longitude < -180 || a.Longitude > 180 {
		return InputError{Param: "area.longitude", Reason: "must be between -180 and 180"}
	}
	return nil
}

func checkSearchInput(in TextSearchInput) error {
	if trimmed(in.Text) == "" {
		return InputError{Param: "text", Reason: "cannot be empty"}
	}
	if in.MaxResults < 1 || in.MaxResults > SearchMaxResults {
		return InputError{Param: "max_results", Reason: fmt.Sprintf("must be 1-%d", SearchMaxResults)}
	}
	if in.Filters != nil {
		if in.Filters.MinScore != nil && (*in.Filters.MinScore < 0 || *in.Filters.MinScore > 5) {
			return InputError{Param: "filters.min_score", Reason: "must be 0-5"}
		}
		for _, tier := range in.Filters.PriceTiers {
			if tier < 0 || tier > 4 {
				return InputError{Param: "filters.price_tiers", Reason: "must be 0-4"}
			}
		}
	}
	if in.Vicinity != nil {
		if err := checkArea(in.Vicinity); err != nil {
			return err
		}
	}
	return nil
}

func checkSuggestInput(in SuggestInput) error {
	if trimmed(in.Fragment) == "" {
		return InputError{Param: "fragment", Reason: "cannot be empty"}
	}
	if in.MaxResults < 1 || in.MaxResults > SuggestMaxResults {
		return InputError{Param: "max_results", Reason: fmt.Sprintf("must be 1-%d", SuggestMaxResults)}
	}
	if in.Vicinity != nil {
		if err := checkArea(in.Vicinity); err != nil {
			return err
		}
	}
	return nil
}

func checkNearbyInput(in NearbyInput) error {
	if in.Center == nil {
		return InputError{Param: "center", Reason: "required"}
	}
	if err := checkArea(in.Center); err != nil {
		return err
	}
	if in.MaxResults < 1 || in.MaxResults > NearbyMaxResults {
		return InputError{Param: "max_results", Reason: fmt.Sprintf("must be 1-%d", NearbyMaxResults)}
	}
	return nil
}

func checkLookupInput(in LocationLookupInput) error {
	if trimmed(in.Address) == "" {
		return InputError{Param: "address", Reason: "cannot be empty"}
	}
	if in.MaxResults < 1 || in.MaxResults > ResolveMaxResults {
		return InputError{Param: "max_results", Reason: fmt.Sprintf("must be 1-%d", ResolveMaxResults)}
	}
	return nil
}
