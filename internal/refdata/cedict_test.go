package refdata_test

import (
	"testing"

	"github.com/marcelocantos/flac/internal/pinyin"
	"github.com/marcelocantos/flac/internal/refdata"
	"github.com/spf13/afero"
)

func BenchmarkNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := refdata.New(
			pinyin.Cache{},
			afero.NewBasePathFs(afero.NewOsFs(), "../../refdata"))
		if err != nil {
			b.Fatal(err)
		}
	}
}
