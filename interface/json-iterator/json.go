package json

import (
	"bytes"
	"io"

	"github.com/a8m/documentdb"
	jsoniter "github.com/json-iterator/go"
)

func init() {
	documentdb.Serialization = documentdb.SerializationDriver{
		EncoderFactory: func(b *bytes.Buffer) documentdb.JSONEncoder {
			return jsoniter.NewEncoder(b)
		},
		DecoderFactory: func(r io.Reader) documentdb.JSONDecoder {
			return jsoniter.NewDecoder(r)
		},
		Marshal:   jsoniter.Marshal,
		Unmarshal: jsoniter.Unmarshal,
	}
}
