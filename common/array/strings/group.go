package strings

func Group(ss []string, size int) [][]string {
	groups := make([][]string, 0, len(ss)/size)
	group := make([]string, 0, size)
	for _, s := range ss {
		group = append(group, s)
		if len(group) == size {
			groups = append(groups, group)
			group = make([]string, 0, size)
		}
	}
	if len(group) > 0 {
		groups = append(groups, group)
	}
	return groups
}
