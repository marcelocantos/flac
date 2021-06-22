package ui

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecoratePhrase(t *testing.T) {
	colorRE := regexp.MustCompile(`<(\w+)\((.*?)\)>`)
	c := func(s string) string {
		return colorRE.ReplaceAllString(s, "[$1::b]$2[-::-]")
	}
	assert.Equal(t,
		c("🆑:个[<purple(gè)>]"),
		DecorateDefinition("CL:個|个[ge4]"))
	assert.Equal(t,
		c("🆑:个[<purple(gè)>],种[<blue(zhǒng)>]"),
		DecorateDefinition("CL:個|个[ge4],種|种[zhong3]"))
	assert.Equal(t,
		c("🆑:门[<green(mén)>],种[<blue(zhǒng)>],项[<purple(xiàng)>]"),
		DecorateDefinition("CL:門|门[men2],種|种[zhong3],項|项[xiang4]"))

	// classifier for ...
	assert.Equal(t,
		c("令 令 [<blue(lǐng)>] /🆑➤ a ream of paper/"),
		DecorateDefinition("令 令 [ling3] /classifier for a ream of paper/"))
	assert.Equal(t,
		c("味 味 [<purple(wèi)>] /taste/smell/(fig.) (noun suffix) feel/"+
			"quality/sense/(TCM) 🆑➤ ingredients of a medicine prescription/"),
		DecorateDefinition(
			"味 味 [wei4] /taste/smell/(fig.) (noun suffix) feel/quality/sense/"+
				"(TCM) classifier for ingredients of a medicine prescription/"))
}
