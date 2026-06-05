package textutil

func Lines(s string) []string {
	if s == "" {
		return nil
	}
	return []string{s}
}
