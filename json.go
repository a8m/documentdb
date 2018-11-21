package documentdb

import (
	"bytes"
	"encoding/json"
	"io"
)

// JsonEncoder describes json encoder
type JsonEncoder interface {
	Encode(val interface{}) error
}

// JsonDecoder describes json decoder
type JsonDecoder interface {
	Decode(obj interface{}) error
}

// Marshal function type
type Marshal func(v interface{}) ([]byte, error)

// Unmarshal function type
type Unmarshal func(data []byte, v interface{}) error

// EncoderFactory describes function that creates json encoder
type EncoderFactory func(*bytes.Buffer) JsonEncoder

// DecoderFactory describes function that creates json decoder
type DecoderFactory func(io.Reader) JsonDecoder

// SerializationDriver struct holds serialization / deserilization providers
type SerializationDriver struct {
	EncoderFactory EncoderFactory
	DecoderFactory DecoderFactory
	Marshal        Marshal
	Unmarshal      Unmarshal
}

// DefaultSerialization holds default stdlib json driver
var DefaultSerialization = SerializationDriver{
	EncoderFactory: func(b *bytes.Buffer) JsonEncoder {
		return json.NewEncoder(b)
	},
	DecoderFactory: func(r io.Reader) JsonDecoder {
		return json.NewDecoder(r)
	},
	Marshal:   json.Marshal,
	Unmarshal: json.Unmarshal,
}

// Serialization holds driver that is actually used
var Serialization = DefaultSerialization
