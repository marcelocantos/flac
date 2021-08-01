package refdata

import (
	"bytes"
	"crypto/rand"
	_ "embed"
	"io"
	"io/ioutil"
	"math/big"
	"strings"

	"github.com/pierrec/lz4"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"

	"github.com/marcelocantos/flac/internal/pkg/proto/refdata_pb"
)

//go:embed refdata.cache
var refdata_proto_lz4 []byte

func New() (*refdata_pb.RefData, error) {
	return NewFromBytes(refdata_proto_lz4)
}

func NewFromBytes(data []byte) (*refdata_pb.RefData, error) {
	var reader io.Reader = bytes.NewBuffer(data)
	if !bytes.HasPrefix(data, []byte("NOCOMPRESS:")) {
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

func RandomDefinition(
	word string,
	entries *refdata_pb.CEDict_Entries,
) (string, *refdata_pb.CEDict_Entries) {
	if len(entries.Entries) == 1 {
		return "", entries
	}

	ret := &refdata_pb.CEDict_Entries{}

	bign, err := rand.Int(rand.Reader, big.NewInt(int64(len(entries.Entries))))
	if err != nil {
		panic(err)
	}
	n := int(bign.Int64())

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
	see := -1
	for _, def := range defs.Definitions {
		switch {
		case strings.HasPrefix(def, "also written "):
		case strings.HasPrefix(def, "also pr. "):
		case strings.HasPrefix(def, "CL:"):
		case (strings.HasPrefix(def, "variant of ") || strings.HasPrefix(def, "see ")) &&
			strings.Contains(def, pinyin):
			candidateDefs = append(candidateDefs,
				strings.ReplaceAll(def, pinyin, "🙈"))
			see = len(candidateDefs)
		case strings.HasPrefix(def, "surname "):
			candidateDefs = append(candidateDefs, "surname")
		default:
			candidateDefs = append(candidateDefs, def)
		}
	}

	// "(see|variant of) ..." aren't great choices of definitions to test. Avoid
	// unless no other options reman.
	if see != -1 && len(candidateDefs) > 1 {
		candidateDefs = append(candidateDefs[:see], candidateDefs[see+1:]...)
	}

	if len(candidateDefs) == 0 {
		panic(errors.Errorf("no useful definitions for %s: %v", word, defs))
	}
	bign, err = rand.Int(rand.Reader, big.NewInt(int64(len(candidateDefs))))
	if err != nil {
		panic(err)
	}
	n = int(bign.Int64())
	def := candidateDefs[n]

	return def, ret
}
