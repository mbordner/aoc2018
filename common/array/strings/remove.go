package strings

func Remove(a []string, s string) []string {
	r := make([]string, 0, len(a))
	for i := range a {
		if s != a[i] {
			r = append(r, a[i])
		}
	}
	return r
}
