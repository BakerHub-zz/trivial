package slices

// Man, we do need generics here!

func Unique(slices []string) []string {
	keys := make(map[string]bool)
	var rv []string
	for _, entry := range slices {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			rv = append(rv, entry)
		}
	}
	return rv
}

func Compact(slices []string) []string {
	var rv []string
	for _, s := range slices {
		if len(s) > 0 {
			rv = append(rv, s)
		}
	}
	return rv
}

func Map(slices []string, fn func(string) string) []string {
	var rv []string
	for _, s := range slices {
		rv = append(rv, fn(s))
	}
	return rv
}

func Filter(slices []string, fn func(string) bool) []string {
	var rv []string
	for _, s := range slices {
		if fn(s) {
			rv = append(rv, s)
		}
	}
	return rv
}

func Contains(slices []string, value string) bool {
	for _, s := range slices {
		if s == value {
			return true
		}
	}
	return false
}
