package refdata

import (
	"bytes"
	_ "embed"
	"io"
	"io/ioutil"
	"math/rand"
	"strings"

	"github.com/pierrec/lz4"
	"google.golang.org/protobuf/proto"

	"github.com/marcelocantos/flac/internal/pkg/proto/refdata_pb"
)

//go:embed refdata.cache
var refdata_proto_lz4 []byte

func New() (*refdata_pb.RefData, error) {
	var reader io.Reader = bytes.NewBuffer(refdata_proto_lz4)
	if !bytes.HasPrefix(refdata_proto_lz4, []byte("NOCOMPRESS:")) {
		reader = lz4.NewReader(reader)
	} else {
		if _, err := reader.Read([]byte("NOCOMPRESS:")); err != nil {
			return nil, err
		}
	}
	refdata_proto, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	rd := &refdata_pb.RefData{}
	err = proto.Unmarshal(refdata_proto, rd)
	return rd, err
}

func RandomDefinition(entries *refdata_pb.CEDict_Entries) (string, *refdata_pb.CEDict_Entries) {
	if len(entries.Entries) == 1 {
		return "", entries
	}

	ret := &refdata_pb.CEDict_Entries{}

	n := rand.Int() % len(entries.Entries)
	var pinyin string
	var defs *refdata_pb.CEDict_Definitions
	for pinyin, defs = range entries.Entries {
		if n == 0 {
			ret.Entries = map[string]*refdata_pb.CEDict_Definitions{pinyin: defs}
			break
		}
		n--
	}

	var candidateDefs []string
	for _, def := range defs.Definitions {
		switch {
		case strings.HasPrefix(def, "surname "):
			candidateDefs = append(candidateDefs, "surname")
		case strings.HasPrefix(def, "also written "):
		case strings.HasPrefix(def, "also pr. "):
		default:
			candidateDefs = append(candidateDefs, def)
		}
	}

	n = rand.Int() % len(candidateDefs)
	def := candidateDefs[n]

	return def, ret
}
