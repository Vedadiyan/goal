package protoutil

func Or(params ...string) string {
	for i := 0; i < len(params); i++ {
		if params[i] != "" {
			return params[i]
		}
	}
	return ""
}
