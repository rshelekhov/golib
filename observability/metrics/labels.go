package metrics

// MapToLabelValues converts map[string]string to []string in the given order of keys
func MapToLabelValues(labels map[string]string, keys []string) []string {
	vals := make([]string, len(keys))
	for i, k := range keys {
		vals[i] = labels[k]
	}
	return vals
}
