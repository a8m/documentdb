package documentdb

import (
)

type Config struct {
	MasterKey	string
}

type DocumentDB struct {
	client	*Client
}

// Create DocumentClient
func New(url string, config Config) *DocumentDB {
	client := &Client{}
	client.Url = url
	client.Config = config
	return &DocumentDB{client}
}

// Read database by self link
func (c *DocumentDB) ReadDatabase(link string) (*Database, error) {
	db := &Database{}
	c.client.Read(link, "dbs", db)
	return nil, nil
}
