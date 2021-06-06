package refdata

// Freq returns the frequency of a word in the list, or -1 if not found.
func (wl *WordList) Freq(word string) int {
	if i, has := wl.Frequencies[word]; has {
		return int(i)
	}
	return -1
}

// Pos returns the position of a word in the list, or -1 if not found.
func (wl *WordList) Pos(word string) int {
	if i, has := wl.Positions[word]; has {
		return int(i)
	}
	return -1
}

// Has returns true iff a word is in the list.
func (wl *WordList) Has(word string) bool {
	return wl.Pos(word) >= 0
}
