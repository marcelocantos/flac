package refdata

import (
	"bytes"
	_ "embed"
	"io/ioutil"

	"github.com/pierrec/lz4"
	"google.golang.org/protobuf/proto"

	"github.com/marcelocantos/flac/internal/pkg/proto/refdata"
	"github.com/marcelocantos/flac/internal/pkg/refdata/words"
)

var (
	//go:embed refdata.cache
	refdata_proto_lz4 []byte
)

type RefData struct {
	*refdata.RefData
}

func New() (RefData, error) {
	refdata_proto, err := ioutil.ReadAll(
		lz4.NewReader(bytes.NewBuffer(refdata_proto_lz4)))
	if err != nil {
		return RefData{}, err
	}

	rd := RefData{RefData: &refdata.RefData{}}
	if err := proto.Unmarshal(refdata_proto, rd); err != nil {
		return RefData{}, err
	}

	return rd, nil
}

func (rd RefData) WordList() words.WordList {
	return words.WordList{WordList: rd.RefData.WordList}
}
