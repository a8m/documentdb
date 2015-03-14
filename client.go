package documentdb

import (
	"encoding/json"
	"strings"
	"net/http"
	"bytes"
	"io"
)

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

	resp, err := c.Do(r.Request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return readJson(resp.Body, ret)
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
