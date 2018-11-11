//
// This project start as a fork of `github.com/nerdylikeme/go-documentdb` version
// but changed, and may be changed later
//
// Goal: add the full functionality of documentdb, align with the other sdks
// and make it more testable
//
package documentdb

import (
	"bytes"
	"net/http"
	"sync"
)

type RequestOptions struct {
	PartitionKey interface{}
}

type Config struct {
	MasterKey  string
	HttpClient http.Client
}

type DocumentDB struct {
	client Clienter
}

// Create DocumentDBClient
func New(url string, config Config) *DocumentDB {
	client := &Client{
		Client: config.HttpClient,
		buffers: &sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer([]byte{})
			},
		},
	}
	client.Url = url
	client.Config = config
	return &DocumentDB{client}
}

// TODO: Add `requestOptions` arguments
// Read database by self link
func (c *DocumentDB) ReadDatabase(link Link, opts ...CallOption) (db *Database, err error) {
	_, err = c.client.Read(link, &db, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// Read collection by self link
func (c *DocumentDB) ReadCollection(link Link, opts ...CallOption) (coll *Collection, err error) {
	_, err = c.client.Read(link, &coll, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// Read document by self link
func (c *DocumentDB) ReadDocument(link Link, doc interface{}, opts ...CallOption) (err error) {
	_, err = c.client.Read(link, &doc, opts...)
	return
}

// Read sporc by self link
func (c *DocumentDB) ReadStoredProcedure(link Link, opts ...CallOption) (sproc *Sproc, err error) {
	_, err = c.client.Read(link, &sproc, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// Read udf by self link
func (c *DocumentDB) ReadUserDefinedFunction(link Link, opts ...CallOption) (udf *UDF, err error) {
	_, err = c.client.Read(link, &udf, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// Read all databases
func (c *DocumentDB) ReadDatabases(opts ...CallOption) (dbs []Database, err error) {
	return c.QueryDatabases(nil, opts...)
}

// Read all collections by db selflink
func (c *DocumentDB) ReadCollections(db Link, opts ...CallOption) (colls []Collection, err error) {
	return c.QueryCollections(db, nil, opts...)
}

// Read all sprocs by collection self link
func (c *DocumentDB) ReadStoredProcedures(coll Link, opts ...CallOption) (sprocs []Sproc, err error) {
	return c.QueryStoredProcedures(coll, nil, opts...)
}

// Read pall udfs by collection self link
func (c *DocumentDB) ReadUserDefinedFunctions(coll Link, opts ...CallOption) (udfs []UDF, err error) {
	return c.QueryUserDefinedFunctions(coll, nil, opts...)
}

// Read all collection documents by self link
// TODO: use iterator for heavy transactions
func (c *DocumentDB) ReadDocuments(coll Link, docs interface{}, opts ...CallOption) (r *Response, err error) {
	return c.QueryDocuments(coll, nil, docs, opts...)
}

// Read all databases that satisfy a query
func (c *DocumentDB) QueryDatabases(query *Query, opts ...CallOption) (dbs Databases, err error) {
	data := struct {
		Databases Databases `json:"Databases,omitempty"`
		Count     int       `json:"_count,omitempty"`
	}{}
	if query != nil {
		_, err = c.client.Query(Link{"dbs"}, query, &data, opts...)
	} else {
		_, err = c.client.Read(Link{"dbs"}, &data, opts...)
	}
	if dbs = data.Databases; err != nil {
		dbs = nil
	}
	return
}

// Read all db-collection that satisfy a query
func (c *DocumentDB) QueryCollections(db Link, query *Query, opts ...CallOption) (colls []Collection, err error) {
	data := struct {
		Collections []Collection `json:"DocumentCollections,omitempty"`
		Count       int          `json:"_count,omitempty"`
	}{}
	if query != nil {
		_, err = c.client.Query(append(db, "colls/"), query, &data, opts...)
	} else {
		_, err = c.client.Read(append(db, "colls/"), &data, opts...)
	}
	if colls = data.Collections; err != nil {
		colls = nil
	}
	return
}

// Read all collection `sprocs` that satisfy a query
func (c *DocumentDB) QueryStoredProcedures(coll Link, query *Query, opts ...CallOption) (sprocs []Sproc, err error) {
	data := struct {
		Sprocs []Sproc `json:"StoredProcedures,omitempty"`
		Count  int     `json:"_count,omitempty"`
	}{}
	if query != nil {
		_, err = c.client.Query(append(coll, "sprocs/"), query, &data, opts...)
	} else {
		_, err = c.client.Read(append(coll, "sprocs/"), &data, opts...)
	}
	if sprocs = data.Sprocs; err != nil {
		sprocs = nil
	}
	return
}

// Read all collection `udfs` that satisfy a query
func (c *DocumentDB) QueryUserDefinedFunctions(coll Link, query *Query, opts ...CallOption) (udfs []UDF, err error) {
	data := struct {
		Udfs  []UDF `json:"UserDefinedFunctions,omitempty"`
		Count int   `json:"_count,omitempty"`
	}{}
	if query != nil {
		_, err = c.client.Query(append(coll, "udfs/"), query, &data, opts...)
	} else {
		_, err = c.client.Read(append(coll, "udfs/"), &data, opts...)
	}
	if udfs = data.Udfs; err != nil {
		udfs = nil
	}
	return
}

// Read all documents in a collection that satisfy a query
func (c *DocumentDB) QueryDocuments(coll Link, query *Query, docs interface{}, opts ...CallOption) (response *Response, err error) {
	data := struct {
		Documents interface{} `json:"Documents,omitempty"`
		Count     int         `json:"_count,omitempty"`
	}{Documents: docs}
	if query != nil {
		response, err = c.client.Query(append(coll, "docs/"), query, &data, opts...)
	} else {
		response, err = c.client.Read(append(coll, "docs/"), &data, opts...)
	}
	return
}

// Create database
func (c *DocumentDB) CreateDatabase(body interface{}, opts ...CallOption) (db *Database, err error) {
	_, err = c.client.Create(Link{"dbs"}, body, &db, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// Create collection
func (c *DocumentDB) CreateCollection(db Link, body interface{}, opts ...CallOption) (coll *Collection, err error) {
	_, err = c.client.Create(append(db, "colls/"), body, &coll, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// Create stored procedure
func (c *DocumentDB) CreateStoredProcedure(coll Link, body interface{}, opts ...CallOption) (sproc *Sproc, err error) {
	_, err = c.client.Create(append(coll, "sprocs/"), body, &sproc, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// Create user defined function
func (c *DocumentDB) CreateUserDefinedFunction(coll Link, body interface{}, opts ...CallOption) (udf *UDF, err error) {
	_, err = c.client.Create(append(coll, "udfs/"), body, &udf, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// Create document
func (c *DocumentDB) CreateDocument(coll Link, doc interface{}, opts ...CallOption) (*Response, error) {
	// id := reflect.ValueOf(doc).Elem().FieldByName("Id")
	// if id.IsValid() && id.String() == "" {
	// 	id.SetString(uuid())
	// }
	return c.client.Create(append(coll, "docs/"), doc, &doc, opts...)
}

// Upsert document
func (c *DocumentDB) UpsertDocument(coll Link, doc interface{}, opts ...CallOption) (*Response, error) {
	// id := reflect.ValueOf(doc).Elem().FieldByName("Id")
	// if id.IsValid() && id.String() == "" {
	// 	id.SetString(uuid())
	// }
	return c.client.Upsert(append(coll, "docs/"), doc, &doc, opts...)
}

// Delete database
func (c *DocumentDB) Delete(link Link) (*Response, error) {
	return c.client.Delete(link)
}

// Replace database
func (c *DocumentDB) ReplaceDatabase(link Link, body interface{}, opts ...CallOption) (db *Database, err error) {
	_, err = c.client.Replace(link, body, &db)
	if err != nil {
		return nil, err
	}
	return
}

// Replace document
func (c *DocumentDB) ReplaceDocument(link Link, doc interface{}, opts ...CallOption) (*Response, error) {
	return c.client.Replace(link, doc, &doc, opts...)
}

// Replace stored procedure
func (c *DocumentDB) ReplaceStoredProcedure(link Link, body interface{}, opts ...CallOption) (sproc *Sproc, err error) {
	_, err = c.client.Replace(link, body, &sproc, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// Replace stored procedure
func (c *DocumentDB) ReplaceUserDefinedFunction(link Link, body interface{}, opts ...CallOption) (udf *UDF, err error) {
	_, err = c.client.Replace(link, body, &udf, opts...)
	if err != nil {
		return nil, err
	}
	return
}

// Execute stored procedure
func (c *DocumentDB) ExecuteStoredProcedure(link Link, params, body interface{}, opts ...CallOption) (err error) {
	_, err = c.client.Execute(link, params, &body, opts...)
	return
}
