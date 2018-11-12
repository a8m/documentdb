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
	HEADER_XDATE             = "X-Ms-Date"
	HEADER_AUTH              = "Authorization"
	HEADER_VER               = "X-Ms-Version"
	HEADER_CONTYPE           = "Content-Type"
	HEADER_CONLEN            = "Content-Length"
	HEADER_IS_QUERY          = "X-Ms-Documentdb-Isquery"
	HEADER_UPSERT            = "x-ms-documentdb-is-upsert"
	HEADER_PARTITION_KEY     = "x-ms-documentdb-partitionkey"
	HEADER_MAX_ITEM_COUNT    = "x-ms-max-item-count"
	HEADER_CONTINUATION      = "x-ms-continuation"
	HEADER_CONSISTENCY       = "x-ms-consistency-level"
	HEADER_SESSION           = "x-ms-session-token"
	HEADER_CROSSPARTITION    = "x-ms-documentdb-query-enablecrosspartitions"
	HEADER_IFMATCH           = "If-Match"
	HEADER_IF_NONE_MATCH     = "If-None-Match"
	HEADER_IF_MODIFIED_SINCE = "If-Modified-Since"
	HEADER_ACTIVITY_ID       = "x-ms-activity-id"
	HEADER_SESSION_TOKEN     = "x-ms-session-token"
	HEADER_REQUEST_CHARGE    = "x-ms-request-charge"

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
func ResourceRequest(link string, req *http.Request) *Request {
	rId, rType := parse(link)
	return &Request{rId, rType, req}
}

// Add 3 default headers to *Request
// "x-ms-date", "x-ms-version", "authorization"
func (req *Request) DefaultHeaders(mKey *Key) (err error) {
	req.Header.Add(HEADER_XDATE, formatDate(time.Now()))
	req.Header.Add(HEADER_VER, SupportedVersion)

	b := buffers.Get().(*bytes.Buffer)
	b.Reset()
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

	buffers.Put(b)

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

func formatDate(t time.Time) string {
	t = t.UTC()
	return t.Format("Mon, 02 Jan 2006 15:04:05 GMT")
}
