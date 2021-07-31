package ui

type keyedStrings [][2]string

func (s keyedStrings) Len() int {
	return len(s)
}

func (s keyedStrings) Less(i, j int) bool {
	return s[i][0] < s[j][0]
}

func (s keyedStrings) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
