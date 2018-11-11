package documentdb

/*
import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ClientStub struct {
	mock.Mock
}

func (c *ClientStub) Read(link string, ret interface{}, opts ...CallOption) (*Response, error) {
	args := c.Called(link)
	return nil, args.Error(0)
}

func (c *ClientStub) Query(link string, query *Query, ret interface{}, opts ...CallOption) (*Response, error) {
	c.Called(link, query)
	return nil, nil
}

func (c *ClientStub) Create(link string, body, ret interface{}, opts ...CallOption) (*Response, error) {
	c.Called(link, body)
	return nil, nil
}

func (c *ClientStub) Upsert(link string, body, ret interface{}, opts ...CallOption) (*Response, error) {
	c.Called(link, body)
	return nil, nil
}

func (c *ClientStub) Delete(link string, opts ...CallOption) (*Response, error) {
	c.Called(link)
	return nil, nil
}

func (c *ClientStub) Replace(link string, body, ret interface{}, opts ...CallOption) (*Response, error) {
	c.Called(link, body)
	return nil, nil
}

func (c *ClientStub) Execute(link string, body, ret interface{}, opts ...CallOption) (*Response, error) {
	c.Called(link, body)
	return nil, nil
}

func TestNew(t *testing.T) {
	assert := assert.New(t)
	client := New("url", Config{MasterKey: "config"})
	assert.IsType(client, &DocumentDB{}, "Should return DocumentDB object")
}

func TestReadDatabaseFailure(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Read", "self_link").Return(errors.New("couldn't read database"))
	db, err := c.ReadDatabase("self_link")
	assert.Nil(t, db)
	assert.EqualError(t, err, "couldn't read database")
}

func TestReadDatabase(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Read", "self_link").Return(nil)
	c.ReadDatabase("self_link")
	client.AssertCalled(t, "Read", "self_link")
}

func TestReadCollection(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Read", "self_link").Return(nil)
	c.ReadCollection("self_link")
	client.AssertCalled(t, "Read", "self_link")
}

func TestReadDocument(t *testing.T) {
	type MyDocument struct {
		Document
		// Your external fields
		Name    string `json:"name,omitempty"`
		Email   string `json:"email,omitempty"`
		IsAdmin bool   `json:"isAdmin,omitempty"`
	}
	var doc MyDocument
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Read", "self_link_doc").Return(nil)
	c.ReadDocument("self_link_doc", &doc)
	client.AssertCalled(t, "Read", "self_link_doc")
}

func TestReadStoredProcedure(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Read", "self_link").Return(nil)
	c.ReadStoredProcedure("self_link")
	client.AssertCalled(t, "Read", "self_link")
}

func TestReadUserDefinedFunction(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Read", "self_link").Return(nil)
	c.ReadUserDefinedFunction("self_link")
	client.AssertCalled(t, "Read", "self_link")
}

func TestReadDatabases(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Read", "dbs").Return(nil)
	c.ReadDatabases()
	client.AssertCalled(t, "Read", "dbs")
}

func TestReadCollections(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	dbLink := "dblink/"
	client.On("Read", dbLink+"colls/").Return(nil)
	c.ReadCollections(dbLink)
	client.AssertCalled(t, "Read", dbLink+"colls/")
}

func TestReadStoredProcedures(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	collLink := "colllink/"
	client.On("Read", collLink+"sprocs/").Return(nil)
	c.ReadStoredProcedures(collLink)
	client.AssertCalled(t, "Read", collLink+"sprocs/")
}

func TestReadUserDefinedFunctions(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	collLink := "colllink/"
	client.On("Read", collLink+"udfs/").Return(nil)
	c.ReadUserDefinedFunctions(collLink)
	client.AssertCalled(t, "Read", collLink+"udfs/")
}

func TestReadDocuments(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	collLink := "colllink/"
	client.On("Read", collLink+"docs/").Return(nil)
	c.ReadDocuments(collLink, struct{}{})
	client.AssertCalled(t, "Read", collLink+"docs/")
}

func TestQueryDatabases(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Query", "dbs", "SELECT * FROM ROOT r").Return(nil)
	c.QueryDatabases(NewQuery("SELECT * FROM ROOT r"))
	client.AssertCalled(t, "Query", "dbs", "SELECT * FROM ROOT r")
}

func TestQueryCollections(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Query", "db_self_link/colls/", "SELECT * FROM ROOT r").Return(nil)
	c.QueryCollections("db_self_link/", NewQuery("SELECT * FROM ROOT r"))
	client.AssertCalled(t, "Query", "db_self_link/colls/", "SELECT * FROM ROOT r")
}

func TestQueryStoredProcedures(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Query", "colls_self_link/sprocs/", "SELECT * FROM ROOT r").Return(nil)
	c.QueryStoredProcedures("colls_self_link/", NewQuery("SELECT * FROM ROOT r"))
	client.AssertCalled(t, "Query", "colls_self_link/sprocs/", "SELECT * FROM ROOT r")
}

func TestQueryUserDefinedFunctions(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Query", "colls_self_link/udfs/", "SELECT * FROM ROOT r").Return(nil)
	c.QueryUserDefinedFunctions("colls_self_link/", NewQuery("SELECT * FROM ROOT r"))
	client.AssertCalled(t, "Query", "colls_self_link/udfs/", "SELECT * FROM ROOT r")
}

func TestQueryDocuments(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	collLink := "coll_self_link/"
	client.On("Query", collLink+"docs/", "SELECT * FROM ROOT r").Return(nil)
	c.QueryDocuments(collLink, NewQuery("SELECT * FROM ROOT r"), struct{}{})
	client.AssertCalled(t, "Query", collLink+"docs/", "SELECT * FROM ROOT r")
}

func TestCreateDatabase(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Create", "dbs", "{}").Return(nil)
	c.CreateDatabase("{}")
	client.AssertCalled(t, "Create", "dbs", "{}")
}

func TestCreateCollection(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Create", "dbs/colls/", "{}").Return(nil)
	c.CreateCollection("dbs/", "{}")
	client.AssertCalled(t, "Create", "dbs/colls/", "{}")
}

func TestCreateStoredProcedure(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Create", "dbs/colls/sprocs/", `{"id":"fn"}`).Return(nil)
	c.CreateStoredProcedure("dbs/colls/", `{"id":"fn"}`)
	client.AssertCalled(t, "Create", "dbs/colls/sprocs/", `{"id":"fn"}`)
}

func TestCreateUserDefinedFunction(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Create", "dbs/colls/udfs/", `{"id":"fn"}`).Return(nil)
	c.CreateUserDefinedFunction("dbs/colls/", `{"id":"fn"}`)
	client.AssertCalled(t, "Create", "dbs/colls/udfs/", `{"id":"fn"}`)
}

func TestCreateDocument(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	// TODO: test error situation, without id, etc...
	var doc Document
	client.On("Create", "dbs/colls/docs/", &doc).Return(nil)
	c.CreateDocument("dbs/colls/", &doc)
	client.AssertCalled(t, "Create", "dbs/colls/docs/", &doc)
	assert.NotEqual(t, doc.Id, "")
}

func TestUpsertDocument(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	// TODO: test error situation, without id, etc...
	var doc Document
	client.On("Upsert", "dbs/colls/docs/", &doc).Return(nil)
	c.UpsertDocument("dbs/colls/", &doc)
	client.AssertCalled(t, "Upsert", "dbs/colls/docs/", &doc)
	assert.NotEqual(t, doc.Id, "")
}

func TestDeleteResource(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}

	client.On("Delete", "self_link_db").Return(nil)
	c.Delete("self_link_db")
	client.AssertCalled(t, "Delete", "self_link_db")
}

func TestReplaceDatabase(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Replace", "db_link", "{}").Return(nil)
	c.ReplaceDatabase("db_link", "{}")
	client.AssertCalled(t, "Replace", "db_link", "{}")
}

func TestReplaceDocument(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Replace", "doc_link", "{}").Return(nil)
	c.ReplaceDocument("doc_link", "{}")
	client.AssertCalled(t, "Replace", "doc_link", "{}")
}

func TestReplaceStoredProcedure(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Replace", "sproc_link", "{}").Return(nil)
	c.ReplaceStoredProcedure("sproc_link", "{}")
	client.AssertCalled(t, "Replace", "sproc_link", "{}")
}

func TestReplaceUserDefinedFunction(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Replace", "udf_link", "{}").Return(nil)
	c.ReplaceUserDefinedFunction("udf_link", "{}")
	client.AssertCalled(t, "Replace", "udf_link", "{}")
}

func TestExecuteStoredProcedure(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client}
	client.On("Execute", "sproc_link", "{}").Return(nil)
	c.ExecuteStoredProcedure("sproc_link", "{}", struct{}{})
	client.AssertCalled(t, "Execute", "sproc_link", "{}")
}
*/
