package json

import (
	"bytes"
	"io"

	"github.com/a8m/documentdb-go"
	jsoniter "github.com/json-iterator/go"
)

func init() {
	documentdb.Serialization = documentdb.SerializationDriver{
		EncoderFactory: func(b *bytes.Buffer) documentdb.JsonEncoder {
			return jsoniter.NewEncoder(b)
		},
		DecoderFactory: func(r io.Reader) documentdb.JsonDecoder {
			return jsoniter.NewDecoder(r)
		},
		Marshal:   jsoniter.Marshal,
		Unmarshal: jsoniter.Unmarshal,
	}
}
