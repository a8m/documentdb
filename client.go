package documentdb

import (
	"encoding/json"
	"strings"
	"net/http"
	"bytes"
	"io"
	"fmt"
)

type Clienter interface {
	Read(link string, ret interface{}) error
	Query(link string, query string, ret interface{}) error
}

type Client struct {
	Url	string
	Config	Config
	http.Client
}

// Read resource by self link
func (c *Client) Read(link string, ret interface{}) error {
	req, err := http.NewRequest("GET", path(c.Url, link), &bytes.Buffer{})
	if err != nil {
		return err
	}

	r := ResourceRequest(link, req)
	if err = r.DefaultHeaders(c.Config.MasterKey); err != nil {
		return err
	}

	return c.do(r, ret)
}

// Query resource
func (c *Client) Query(link, query string, ret interface{}) error {
	buf := bytes.NewBufferString(querify(query))
	req, err := http.NewRequest("POST", path(c.Url, link), buf)
	if err != nil {
		return err
	}
	r := ResourceRequest(link, req)
	if err = r.DefaultHeaders(c.Config.MasterKey); err != nil {
		return err
	}
	r.QueryHeaders(buf.Len())
	return c.do(r, ret)
}

// Private Do function, DRY
func (c *Client) do(r *Request, data interface{}) error {
	resp, err := c.Do(r.Request)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		err = &RequestError{}
		readJson(resp.Body, &err)
		return err
	}
	defer resp.Body.Close()
	return readJson(resp.Body, data)
}

// Generate link
func path(url string, args ...string) (link string) {
	args = append([]string{url}, args...)
	link = strings.Join(args, "/")
	return
}

// Read json response to given interface(struct, map, ..)
func readJson(reader io.Reader, data interface{}) error {
	return json.NewDecoder(reader).Decode(&data)
}

// Stringify query-string as documentdb expected
func querify(query string) string {
	return fmt.Sprintf(`{ "%s": "%s" }`, "query", query)
}
