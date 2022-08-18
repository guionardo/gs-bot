package internal

func JoinNotEmpty(s ...string) (r string) {
	r = ""
	for _, v := range s {
		if len(v) == 0 {
			continue
		}
		if len(r) > 0 {
			r += ", "
		}
		r += v
	}
	return
}
