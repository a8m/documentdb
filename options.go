package documentdb

import (
	"encoding/json"
	"strconv"
)

type Consistency string

const (
	Strong   Consistency = "strong"
	Bounded  Consistency = "bounded"
	Session  Consistency = "session"
	Eventual Consistency = "eventual"
)

type CallOption func(r *Request) error

func PartitionKey(partitionKey interface{}) CallOption {

	// The partition key header must be an array following the spec:
	// https: //docs.microsoft.com/en-us/rest/api/cosmos-db/common-cosmosdb-rest-request-headers
	// and must contain brackets
	// example: x-ms-documentdb-partitionkey: [ "abc" ]
	var (
		pk  []byte
		err error
	)
	switch v := partitionKey.(type) {
	case json.Marshaler:
		pk, err = json.Marshal(v)
	default:
		pk, err = json.Marshal([]interface{}{v})
	}

	return func(r *Request) error {
		if err != nil {
			return err
		}
		r.Header.Set(HEADER_PARTITION_KEY, string(pk))
		return nil
	}
}

func Upsert() CallOption {
	return func(r *Request) error {
		r.Header.Set(HEADER_UPSERT, "true")
		return nil
	}
}

func Limit(limit int) CallOption {
	header := strconv.Itoa(limit)
	return func(r *Request) error {
		r.Header.Set(HEADER_MAX_ITEM_COUNT, header)
		return nil
	}
}

func Continuation(continuation string) CallOption {
	return func(r *Request) error {
		if continuation == "" {
			return nil
		}
		r.Header.Set(HEADER_CONTINUATION, continuation)
		return nil
	}
}

func ConsistencyLevel(consistency Consistency) CallOption {
	return func(r *Request) error {
		r.Header.Set(HEADER_CONSISTENCY, string(consistency))
		return nil
	}
}

func SessionToken(sessionToken string) CallOption {
	return func(r *Request) error {
		r.Header.Set(HEADER_SESSION, sessionToken)
		return nil
	}
}

func CrossPartition() CallOption {
	return func(r *Request) error {
		r.Header.Set(HEADER_CROSSPARTITION, "True")
		return nil
	}
}
