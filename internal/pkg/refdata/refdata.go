package refdata

import (
	"bytes"
	_ "embed"
	"io/ioutil"

	"github.com/pierrec/lz4"
	"google.golang.org/protobuf/proto"

	"github.com/marcelocantos/flac/internal/pkg/proto/refdata"
)

var (
	//go:embed refdata.cache
	refdata_proto_lz4 []byte
)

func New() (*refdata.RefData, error) {
	refdata_proto, err := ioutil.ReadAll(
		lz4.NewReader(bytes.NewBuffer(refdata_proto_lz4)))
	if err != nil {
		return nil, err
	}

	rd := &refdata.RefData{}
	err = proto.Unmarshal(refdata_proto, rd)
	return rd, err
}
