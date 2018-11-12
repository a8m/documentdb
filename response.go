package documentdb

import "net/http"

type Response struct {
	Header http.Header
}

// Continuation returns continuation token for paged request.
// Pass this value to next request to get next page of documents.
func (r *Response) Continuation() string {
	return r.Header.Get(HEADER_CONTINUATION)
}
