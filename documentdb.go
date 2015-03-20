//
// This project start as a fork of `github.com/nerdylikeme/go-documentdb` version
// but changed a bit, and may be changed later
//
// Goal: add the full functionality of documentdb, align with the other sdks
// and make it more testable
//
package documentdb

type Config struct {
	MasterKey	string
}

type DocumentDB struct {
	client	Clienter
}

// Create DocumentDBClient
func New(url string, config Config) *DocumentDB {
	client := &Client{}
	client.Url = url
	client.Config = config
	return &DocumentDB{client}
}

// TODO: Add `requestOptions` arguments
// Read database by self link
func (c *DocumentDB) ReadDatabase(link string) (db *Database, err error) {
	err = c.client.Read(link, &db)
	if err != nil {
		return nil, err
	}
	return
}

// Read collection by self link
func (c *DocumentDB) ReadCollection(link string) (coll *Collection, err error) {
	err = c.client.Read(link, &coll)
	if err != nil {
		return nil, err
	}
	return
}

// Read document by self link
func (c *DocumentDB) ReadDocument(link string, doc interface{}) (err error) {
	err = c.client.Read(link, &doc)
	return
}

// Read sporc by self link
func (c *DocumentDB) ReadStoredProcedure(link string) (sproc *Sproc, err error) {
	err = c.client.Read(link, &sproc)
	if err != nil {
		return nil, err
	}
	return
}

// Read udf by self link
func (c *DocumentDB) ReadUserDefinedFunction(link string) (udf *UDF, err error) {
	err = c.client.Read(link, &udf)
	if err != nil {
		return nil, err
	}
	return
}

// Read all databases
func (c *DocumentDB) ReadDatabases() (dbs []Database, err error) {
	return c.QueryDatabases("")
}

// Read all collections by db selflink
func (c *DocumentDB) ReadCollections(db string) (colls []Collection, err error) {
	return c.QueryCollections(db, "")
}

// Read all sprocs by collection self link
func (c *DocumentDB) ReadStoredProcedures(coll string) (sprocs []Sproc, err error) {
	data := struct {
		Sprocs	[]Sproc	`json:"StoredProcedures,omitempty"`
		Count	int	`json:"_count,omitempty"`
	}{}
	err = c.client.Read(coll + "sprocs/", &data)
	if sprocs = data.Sprocs; err != nil {
		sprocs = nil
	}
	return
}

// Read all udfs by collection self link
func (c *DocumentDB) ReadUserDefinedFunctions(coll string) (udfs []UDF, err error) {
	data := struct {
		Udfs	[]UDF	`json:"UserDefinedFunctions,omitempty"`
		Count	int	`json:"_count,omitempty"`
	}{}
	err = c.client.Read(coll + "udfs/", &data)
	if udfs = data.Udfs; err != nil {
		udfs = nil
	}
	return
}

// Read all documents by collection self link
// TODO: use read/query iterator
func (c *DocumentDB) ReadDocuments(coll string, docs interface{}) (err error) {
	data := struct {
		Documents	interface{}	`json:"Documents,omitempty"`
		Count		int		`json:"_count,omitempty"`
	}{Documents: docs}
	err = c.client.Read(coll + "docs/", &data)
	return
}

// Read all databases that satisfy a query
func (c *DocumentDB) QueryDatabases(query string) (dbs []Database, err error) {
	data := struct {
		Databases	[]Database	`json:"Databases,omitempty"`
		Count		int		`json:"_count,omitempty"`
	}{}
	if len(query) > 0 {
		err = c.client.Query("dbs", query, &data)
	} else {
		err = c.client.Read("dbs", &data)
	}
	if dbs = data.Databases; err != nil {
		dbs = nil
	}
	return
}

// Read all db-collection that satisfy a query
func (c *DocumentDB) QueryCollections(db, query string) (colls []Collection, err error) {
	data := struct {
		Collections	[]Collection	`json:"DocumentCollections,omitempty"`
		Count		int		`json:"_count,omitempty"`
	}{}
	if len(query) > 0 {
		err = c.client.Query(db + "colls/", query, &data)
	} else {
		err = c.client.Read(db + "colls/", &data)
	}
	if colls = data.Collections; err != nil {
		colls = nil
	}
	return
}
