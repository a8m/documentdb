package documentdb

import (
	"strings"
	"time"
	"net/http"
	"net/url"
)

const (
	HEADER_XDATE	= "X-Ms-Date"
	HEADER_AUTH 	= "Authorization"
	HEADER_VER	= "X-Ms-Version"
)

type Request struct {
	rId	string
	rType	string
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
	req.Header.Add(HEADER_VER, "2014-08-21")

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


