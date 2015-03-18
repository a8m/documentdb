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
func (c *DocumentDB) ReadDocument(link string, doc interface{}) error {
	err := c.client.Read(link, &doc)
	return err
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
	data := struct {
		Databases	[]Database	`json:"Databases,omitempty"`
		Count		int		`json:"_count,omitempty"`
	}{}
	err = c.client.Read("dbs", &data)
	if err != nil {
		dbs = nil
	} else {
		dbs = data.Databases
	}
	return
}

// Read all collections by db selflink
func (c *DocumentDB) ReadCollections(db string) (colls []Collection, err error) {
	data := struct {
		Collections	[]Collection	`json:"DocumentCollections,omitempty"`
		Count		int		`json:"_count,omitempty"`
	}{}
	err = c.client.Read(db + "colls/", &data)
	if err != nil {
		colls = nil
	} else {
		colls = data.Collections
	}
	return
}

func (c *DocumentDB) ReadStoredProcedures(coll string) (sprocs []Sproc, err error) {
	data := struct {
		Sprocs	[]Sproc	`json:"StoredProcedures,omitempty"`
		Count	int	`json:"_count,omitempty"`
	}{}
	err = c.client.Read(coll + "sprocs/", &data)
	if err != nil {
		sprocs = nil
	} else {
		sprocs = data.Sprocs
	}
	return
}
