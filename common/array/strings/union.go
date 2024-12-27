package strings

func Union(v1, v2 []string) []string {
	vs := make(map[string]bool)
	for _, v := range v1 {
		vs[v] = true
	}
	for _, v := range v2 {
		vs[v] = true
	}
	values := make([]string, 0, len(vs))
	for v := range vs {
		values = append(values, v)
	}
	return values
}
