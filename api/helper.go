package api

func eqIntSclice(a []int, b []int) bool {
	if len(a) != len(b) {
		return false
	}

	for i, n := range a {
		if n != b[i] {
			return false
		}
	}

	return true
}
