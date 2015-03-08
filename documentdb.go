package documentdb

import (
	"fmt"
)

type Config struct {
	MasterKey	string
}

type DocumentDB struct {
	Url	string
	Config	Config
	client	*Client
}

// Create DocumentClient
func New(url string, config Config) *DocumentDB {
	return &DocumentDB{url, config, &Client{}}
}

// Read database by self link
//func ( Client) ReadDatabase(link string) (*Database, error) {
//	link = "/" + link + "/"
//
//}
