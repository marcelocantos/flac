package ui

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccentPhrase(t *testing.T) {
	colorRE := regexp.MustCompile(`<(\w+)\((.*?)\)>`)
	c := func(s string) string {
		return colorRE.ReplaceAllString(s, "[$1::b]$2[-::-]")
	}
	assert.Equal(t,
		c("CL:个[<purple(gè)>]"),
		accentPhrase("CL:個|个[ge4]"))
	assert.Equal(t,
		c("CL:个[<purple(gè)>],种[<blue(zhǒng)>]"),
		accentPhrase("CL:個|个[ge4],種|种[zhong3]"))
	assert.Equal(t,
		c("CL:门[<green(mén)>],种[<blue(zhǒng)>],项[<purple(xiàng)>]"),
		accentPhrase("CL:門|门[men2],種|种[zhong3],項|项[xiang4]"))
}
