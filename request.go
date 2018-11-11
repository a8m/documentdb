package documentdb

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	HEADER_XDATE          = "X-Ms-Date"
	HEADER_AUTH           = "Authorization"
	HEADER_VER            = "X-Ms-Version"
	HEADER_CONTYPE        = "Content-Type"
	HEADER_CONLEN         = "Content-Length"
	HEADER_IS_QUERY       = "X-Ms-Documentdb-Isquery"
	HEADER_UPSERT         = "x-ms-documentdb-is-upsert"
	HEADER_PARTITION_KEY  = "x-ms-documentdb-partitionkey"
	HEADER_MAX_ITEM_COUNT = "x-ms-max-item-count"
	HEADER_CONTINUATION   = "x-ms-continuation"
	HEADER_CONSISTENCY    = "x-ms-consistency-level"
	HEADER_SESSION    = "x-ms-session-token"
	HEADER_CROSSPARTITION    = "x-ms-documentdb-query-enablecrosspartitions"

	SupportedVersion = "2017-02-22"
)

// Request Error
type RequestError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Implement Error function
func (e RequestError) Error() string {
	return fmt.Sprintf("%v, %v", e.Code, e.Message)
}

// Resource Request
type Request struct {
	rId, rType string
	*http.Request
}

// Return new resource request with type and id
func ResourceRequest(link Link, req *http.Request) *Request {
	rId, rType := parse(req.URL.Path)
	return &Request{rId, rType, req}
}

// Add 3 default headers to *Request
// "x-ms-date", "x-ms-version", "authorization"
func (req *Request) DefaultHeaders(mKey string) (err error) {
	req.Header.Add(HEADER_XDATE, time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"))
	req.Header.Add(HEADER_VER, SupportedVersion)

	b := bytes.Buffer{}
	b.WriteString(req.Method)
	b.WriteRune('\n')
	b.WriteString(req.rType)
	b.WriteRune('\n')
	b.WriteString(req.rId)
	b.WriteRune('\n')
	b.WriteString(req.Header.Get(HEADER_XDATE))
	b.WriteRune('\n')
	b.WriteString(req.Header.Get("Date"))
	b.WriteRune('\n')

	sign, err := authorize(bytes.ToLower(b.Bytes()), mKey)
	if err != nil {
		return err
	}

	req.Header.Add(HEADER_AUTH, url.QueryEscape("type=master&ver=1.0&sig="+sign))

	return
}

// Add headers for query request
func (req *Request) QueryHeaders(len int) {
	req.Header.Add(HEADER_CONTYPE, "application/query+json")
	req.Header.Add(HEADER_IS_QUERY, "true")
	req.Header.Add(HEADER_CONLEN, string(len))
}

func parse(id string) (rId, rType string) {
	if strings.HasPrefix(id, "/") == false {
		id = "/" + id
	}
	if strings.HasSuffix(id, "/") == false {
		id = id + "/"
	}

	parts := strings.Split(id, "/")
	l := len(parts) // 4

	if l%2 == 0 {
		rId = parts[l-2]
		rType = parts[l-3]
	} else {
		rId = parts[l-3]
		rType = parts[l-2]
	}
	return
}

// // Get path and return resource Id and Type
// // (e.g: "/dbs/b5NCAA==/" ==> "b5NCAA==", "dbs")
// func parse(link Link) (rId, rType string) {

// 	fmt.Println("parse:", link)

// 	l := len(link) // 4

// 	if l == 1 {
// 		rType = link[0]
// 	} else {

// 		if l%2 == 0 {
// 			rId = link[l-1]
// 			rType = link[l-2]
// 		} else {
// 			rId = link[l-2]
// 			rType = link[l-1]
// 		}
// 	}
// 	l = len(rId)
// 	if l > 0 && rId[l-1] == '/' {
// 		rId = rId[0 : l-1]
// 	}
// 	l = len(rType)
// 	if rType[l-1] == '/' {
// 		rType = rType[0 : l-1]
// 	}
// 	if rType[0] == '/' {
// 		rType = rType[1:]
// 	}
// 	return
// }
