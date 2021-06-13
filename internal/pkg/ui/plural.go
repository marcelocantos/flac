package ui

func plural(n int, singular, plural string) string {
	if n == 1 {
		return singular
	}
	return plural
}

func pluralS(n int) string {
	return plural(n, "", "s")
}
