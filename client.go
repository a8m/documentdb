package documentdb

import (
	"fmt"
	"strings"
	"net/http"
	"bytes"
	"io/ioutil"
)

type Client struct {
	Url	string
	Config	Config
	http.Client
}

func (c *Client) Read(rId, rType string, ret interface{}) error {
	req, err := http.NewRequest("GET", path(c.Url, rId), &bytes.Buffer{})
	if err != nil {
		return err
	}
	r := ResourceRequest(rId, rType, req)
	if err = r.DefaultHeaders(c.Config.MasterKey); err != nil {
		return err
	}
	resp, err := c.Do(r.Request)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	return nil
}

// Generate link
func path(url string, args ...string) (link string) {
	args = append([]string{url}, args...)
	link = strings.Join(args, "/")
	return
}


// Default Headers
//{
// 	'Cache-Control': 'no-cache',
//	'x-ms-version': '2014-08-21',
//	'User-Agent': 'documentdb-nodejs-sdk-0.9.1',
//	'x-ms-date': 'Sat, 07 Mar 2015 20:42:27 GMT',
//	authorization: 'type%3Dmaster%26ver%3D1.0%26sig%3DObTpHah0GgBLJ1KMGITRMM9G5%2F4YRodrBrInUR3t%2B00%3D',
//	Accept: 'application/json'
// }
