package documentdb

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	HeaderXDate               = "X-Ms-Date"
	HeaderAuth                = "Authorization"
	HeaderVersion             = "X-Ms-Version"
	HeaderContentType         = "Content-Type"
	HeaderContentLength       = "Content-Length"
	HeaderIsQuery             = "X-Ms-Documentdb-Isquery"
	HeaderUpsert              = "x-ms-documentdb-is-upsert"
	HeaderPartitionKey        = "x-ms-documentdb-partitionkey"
	HeaderMaxItemCount        = "x-ms-max-item-count"
	HeaderContinuation        = "x-ms-continuation"
	HeaderConsistency         = "x-ms-consistency-level"
	HeaderSessionToken        = "x-ms-session-token"
	HeaderCrossPartition      = "x-ms-documentdb-query-enablecrosspartition"
	HeaderIfMatch             = "If-Match"
	HeaderIfNonMatch          = "If-None-Match"
	HeaderIfModifiedSince     = "If-Modified-Since"
	HeaderActivityID          = "x-ms-activity-id"
	HeaderRequestCharge       = "x-ms-request-charge"
	HeaderAIM                 = "A-IM"
	HeaderPartitionKeyRangeID = "x-ms-documentdb-partitionkeyrangeid"
	HeaderUserAgent           = "User-Agent"

	SupportedVersion = "2017-02-22"

	ServicePrincipalRefreshTimeout = 10 * time.Second
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
func (req *Request) DefaultHeaders(config *Config, userAgent string) (err error) {
	req.Header.Add(HeaderXDate, formatDate(time.Now()))
	req.Header.Add(HeaderVersion, SupportedVersion)
	req.Header.Add(HeaderUserAgent, userAgent)

	// Authentication via master key
	if config.MasterKey != nil && config.MasterKey.Key != "" {
		b := buffers.Get().(*bytes.Buffer)
		b.Reset()
		b.WriteString(strings.ToLower(req.Method))
		b.WriteRune('\n')
		b.WriteString(strings.ToLower(req.rType))
		b.WriteRune('\n')
		b.WriteString(req.rId)
		b.WriteRune('\n')
		b.WriteString(strings.ToLower(req.Header.Get(HeaderXDate)))
		b.WriteRune('\n')
		b.WriteString(strings.ToLower(req.Header.Get("Date")))
		b.WriteRune('\n')

		sign, err := authorize(b.Bytes(), config.MasterKey)
		if err != nil {
			return err
		}

		buffers.Put(b)

		req.Header.Add(HeaderAuth, url.QueryEscape("type=master&ver=1.0&sig="+sign))
	} else if config.ServicePrincipal != nil {
		ctx, cancel := context.WithTimeout(req.Context(), ServicePrincipalRefreshTimeout)
		defer cancel()
		err := config.ServicePrincipal.EnsureFreshWithContext(ctx)
		if err != nil {
			return err
		}
		token := config.ServicePrincipal.OAuthToken()
		req.Header.Add(HeaderAuth, url.QueryEscape("type=aad&ver=1.0&sig="+token))
	}

	return
}

// Add headers for query request
func (req *Request) QueryHeaders(len int) {
	req.Header.Add(HeaderContentType, "application/query+json")
	req.Header.Add(HeaderIsQuery, "true")
	req.Header.Add(HeaderContentLength, strconv.Itoa(len))
}

func parse(id string) (rId, rType string) {
	if strings.HasPrefix(id, "/") == false {
		id = "/" + id
	}
	if strings.HasSuffix(id, "/") == false {
		id = id + "/"
	}

	parts := strings.Split(id, "/")
	l := len(parts)

	if l%2 == 0 {
		rType = parts[l-3]
	} else {
		rType = parts[l-2]
	}

	// Check if we're being passed a _self link or a link that uses IDs
	// If we have a self link, parts[2] should be a 6-byte, base64-encoded string, that is 8 characters long and includes padding ("==")
	// "=" is not a valid character in a Cosmos DB identifier, so if we notice that (especially in a string that's 8-chars long), we know it's a RID
	if l > 3 && len(parts[2]) == 8 && parts[2][6:] == "==" {
		// We have a _self link
		if l%2 == 0 {
			rId = parts[l-2]
		} else {
			rId = parts[l-3]
		}
	} else {
		// We have a link that uses IDs
		end := l - 1
		if l%2 == 1 {
			end = l - 2
		}
		rId = strings.Join(parts[1:end], "/")
	}

	return
}

func formatDate(t time.Time) string {
	t = t.UTC()
	return t.Format("Mon, 02 Jan 2006 15:04:05 GMT")
}

type queryPartitionKeyRangesRequest struct {
	Ranges []PartitionKeyRange `json:"PartitionKeyRanges,omitempty"`
	Count  int                 `json:"_count,omitempty"`
}
