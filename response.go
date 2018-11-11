package documentdb

import "net/http"

type Response struct {
	Header http.Header
}

func (r *Response) Continuation() string {
	return r.Header.Get(HEADER_CONTINUATION)
}
