package documentdb

import (
	"fmt"
	"strings"
	"time"
	"net/http"
	"net/url"
)

const (
	HEADER_XDATE	= "X-Ms-Date"
	HEADER_AUTH 	= "Authorization"
	HEADER_VER	= "X-Ms-Version"
	HEADER_CONTYPE	= "Content-Type"
	HEADER_CONLEN	= "Content-Length"
	HEADER_IS_QUERY	= "X-Ms-Documentdb-Isquery"
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
	rId, rType	string
	*http.Request
}

// Return new resource request with type and id
func ResourceRequest(link string, req *http.Request) *Request {
	rId, rType := parse(link)
	return &Request{rId, rType, req}
}

// Add 3 default headers to *Request
// "x-ms-date", "x-ms-version", "authorization"
func (req *Request) DefaultHeaders(mKey string) (err error) {
	req.Header.Add(HEADER_XDATE, time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"))
	req.Header.Add(HEADER_VER, "2017-02-22")

	// Auth
	parts := []string{req.Method, req.rType, req.rId, req.Header.Get(HEADER_XDATE), req.Header.Get("Date"), ""}
	sign, err := authorize(strings.ToLower(strings.Join(parts, "\n")), mKey)
	if err != nil {
		return err
	}

	masterToken := "master"
	tokenVersion := "1.0"
	req.Header.Add(HEADER_AUTH, url.QueryEscape("type=" + masterToken + "&ver=" + tokenVersion + "&sig=" +sign))
	return
}

// Add headers for query request
func (req *Request) QueryHeaders(len int) {
	req.Header.Add(HEADER_CONTYPE, "application/query+json")
	req.Header.Add(HEADER_IS_QUERY, "true")
	req.Header.Add(HEADER_CONLEN, string(len))
}

// Get path and return resource Id and Type
// (e.g: "/dbs/b5NCAA==/" ==> "b5NCAA==", "dbs")
func parse(id string) (rId, rType string) {
	if strings.HasPrefix(id, "/") == false {
		id = "/" + id
	}
	if strings.HasSuffix(id, "/") == false {
		id = id + "/"
	}

	parts := strings.Split(id, "/")
	l := len(parts)

	if l % 2 == 0 {
		rId = parts[l - 2]
		rType = parts[l - 3]
	} else {
		rId = parts[l - 3]
		rType = parts[l - 2]
	}
	return
}


