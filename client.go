package documentdb

import (
	"bytes"
	"io"
	"net/http"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

var (
	newEncoder = jsoniter.NewEncoder
	newDecoder = jsoniter.NewDecoder
	//newEncoder = json.NewEncoder
	//newDecoder = json.NewDecoder
)

type Clienter interface {
	Read(link Link, ret interface{}, opts ...CallOption) (*Response, error)
	Delete(link Link, opts ...CallOption) (*Response, error)
	Query(link Link, query *Query, ret interface{}, opts ...CallOption) (*Response, error)
	Create(link Link, body, ret interface{}, opts ...CallOption) (*Response, error)
	Upsert(link Link, body, ret interface{}, opts ...CallOption) (*Response, error)
	Replace(link Link, body, ret interface{}, opts ...CallOption) (*Response, error)
	Execute(link Link, body, ret interface{}, opts ...CallOption) (*Response, error)
}

type Client struct {
	Url    string
	Config Config
	http.Client
	buffers *sync.Pool
}

func (c *Client) apply(r *Request, opts []CallOption) (err error) {
	if err = r.DefaultHeaders(c.Config.MasterKey); err != nil {
		return err
	}

	for i := 0; i < len(opts); i++ {
		if err = opts[i](r); err != nil {
			return err
		}
	}
	return nil
}

// Read resource by self link
func (c *Client) Read(link Link, ret interface{}, opts ...CallOption) (*Response, error) {
	buf := c.buffers.Get().(*bytes.Buffer)

	res, err := c.method("GET", link, http.StatusOK, ret, buf, opts...)

	c.buffers.Put(buf)

	return res, err
}

// Delete resource by self link
func (c *Client) Delete(link Link, opts ...CallOption) (*Response, error) {
	return c.method("DELETE", link, http.StatusNoContent, nil, &bytes.Buffer{}, opts...)
}

// Query resource
func (c *Client) Query(link Link, query *Query, ret interface{}, opts ...CallOption) (*Response, error) {
	var (
		err error
		req *http.Request
		buf = c.buffers.Get().(*bytes.Buffer)
	)
	defer func() { buf.Reset(); c.buffers.Put(buf) }()

	if err = newEncoder(buf).Encode(query); err != nil {
		return nil, err

	}

	req, err = http.NewRequest("POST", link.ToURL(c.Url), buf)
	if err != nil {
		return nil, err
	}
	r := ResourceRequest(link, req)

	if err = c.apply(r, opts); err != nil {
		return nil, err
	}

	r.QueryHeaders(buf.Len())

	return c.do(r, http.StatusOK, ret)
}

// Create resource
func (c *Client) Create(link Link, body, ret interface{}, opts ...CallOption) (*Response, error) {
	data, err := stringify(body)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	return c.method("POST", link, http.StatusCreated, ret, buf, opts...)
}

// Upsert resource
func (c *Client) Upsert(link Link, body, ret interface{}, opts ...CallOption) (*Response, error) {
	opts = append(opts, Upsert())
	return c.Create(link, body, ret, opts...)
}

// Replace resource
func (c *Client) Replace(link Link, body, ret interface{}, opts ...CallOption) (*Response, error) {
	data, err := stringify(body)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	return c.method("PUT", link, http.StatusOK, ret, buf, opts...)
}

// Replace resource
// TODO: DRY, move to methods instead of actions(POST, PUT, ...)
func (c *Client) Execute(link Link, body, ret interface{}, opts ...CallOption) (*Response, error) {
	data, err := stringify(body)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	return c.method("POST", link, http.StatusOK, ret, buf, opts...)
}

// Private generic method resource
func (c *Client) method(method string, link Link, status int, ret interface{}, body *bytes.Buffer, opts ...CallOption) (*Response, error) {
	req, err := http.NewRequest(method, link.ToURL(c.Url), body)
	if err != nil {
		return nil, err
	}

	r := ResourceRequest(link, req)

	if err = c.apply(r, opts); err != nil {
		return nil, err
	}

	return c.do(r, status, ret)
}

// Private Do function, DRY
func (c *Client) do(r *Request, status int, data interface{}) (*Response, error) {
	resp, err := c.Do(r.Request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != status {
		err = &RequestError{}
		readJson(resp.Body, &err)
		return nil, err
	}
	defer resp.Body.Close()
	if data == nil {
		return nil, nil
	}
	return &Response{resp.Header}, readJson(resp.Body, data)
}

// Read json response to given interface(struct, map, ..)
func readJson(reader io.Reader, data interface{}) error {
	return newDecoder(reader).Decode(&data)
}

// Stringify body data
func stringify(body interface{}) (bt []byte, err error) {
	switch t := body.(type) {
	case string:
		bt = []byte(t)
	case []byte:
		bt = t
	default:
		bt, err = jsoniter.Marshal(t)
	}
	return
}
