package pinyin

import "math/bits"

type Tones int

func newTonesFromString(s string) Tones {
	tones := 0
	for _, r := range s {
		tones |= 1 << (r - '0')
	}
	return Tones(tones)
}

func (t Tones) Count() int {
	return bits.OnesCount(uint(t))
}

func (t Tones) Range() TonesRanger {
	return TonesRanger(t | 1)
}

type TonesRanger Tones

func (r *TonesRanger) Next() bool {
	*r = *r & (*r - 1)
	return *r != 0
}

func (r *TonesRanger) Value() int {
	return bits.TrailingZeros(uint(*r))
}
