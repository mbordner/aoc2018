package strings

func Intersect(v1, v2 []string) (values []string, extras []string) {
	vs := make(map[string]bool)
	values = make([]string, 0, len(v1)+len(v2))
	extras = make([]string, 0, 5)
	for _, v := range v1 {
		vs[v] = false
	}
	for _, v := range v2 {
		if _, ok := vs[v]; ok {
			vs[v] = true
			values = append(values, v)
		} else {
			extras = append(extras, v)
		}
	}
	for v, ok := range vs {
		if !ok {
			extras = append(extras, v)
		}
	}
	return
}
