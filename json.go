package documentdb

import (
	jsoniter "github.com/json-iterator/go"
)

var (
	newEncoder = jsoniter.NewEncoder
	newDecoder = jsoniter.NewDecoder
	marshal    = jsoniter.Marshal
	unmarshal  = jsoniter.Unmarshal
	//newEncoder = json.NewEncoder
	//newDecoder = json.NewDecoder
)
