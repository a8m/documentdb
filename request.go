package documentdb

import (
	"strings"
	"time"
	"net/http"
	"net/url"
)

const (
	HEADER_XDATE	= "x-ms-date"
	HEADER_AUTH 	= "authorization"
)

type Request struct {
	rId	string
	rType	string
	*http.Request
}

func ResourceRequest(rId, rType string, req *http.Request) *Request {
	return &Request{rId, rType, req}
}

func (req *Request) DefaultHeaders(mKey string) error {
	req.Header.Add(HEADER_XDATE, time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"))
	req.Header.Add("x-ms-version", "2014-08-21")

	// Auth
	parts := []string{req.Method, req.rType, req.rId, req.Header.Get(HEADER_XDATE), req.Header.Get("Date"), ""}
	sign, err := authorize(strings.ToLower(strings.Join(parts, "\n")), mKey)
	if err != nil {
		return err
	}

	masterToken := "master"
	tokenVersion := "1.0"
	req.Header.Add(HEADER_AUTH, url.QueryEscape("type=" + masterToken + "&ver=" + tokenVersion + "&sig=" +sign))
	return nil
}

// Response Example
// { type: 'dbs',
// objectBody: { id: 'b5NCAA==', self: '/dbs/b5NCAA==/' } }
//
// { type: 'colls',
// objectBody: { id: 'b5NCAIu9NwA=', self: '/dbs/b5NCAA==/colls/b5NCAIu9NwA=/' } }
//
// Should return a type/id struct
func parse(id string) {
	if strings.HasPrefix(id, "/") == false {
		id = "/" + id
	}
	if strings.HasSuffix(id, "/") == false {
		id = id + "/"
	}
}


