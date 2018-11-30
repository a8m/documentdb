package documentdb

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
	args := c.Called(link, ret, opts)
	r := args.Get(0)
	if r == nil {
		return nil, args.Error(1)
	}
	return r.(*Response), args.Error(1)
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

var defaultConfig = &Config{
	IdentificationHydrator:     DefaultIdentificationHydrator,
	IdentificationPropertyName: "Id",
}

func TestNew(t *testing.T) {
	assert := assert.New(t)
	client := New("url", NewConfig(&Key{Key: "YXJpZWwNCg=="}))
	assert.IsType(client, &DocumentDB{}, "Should return DocumentDB object")
}

func TestReadDatabaseFailure(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	client.On("Read", "self_link", mock.Anything, mock.Anything).Return(nil, errors.New("couldn't read database"))
	db, err := c.ReadDatabase("self_link")
	assert.Nil(t, db)
	assert.EqualError(t, err, "couldn't read database")
}

func TestReadDatabase(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	client.On("Read", "self_link", mock.Anything, mock.Anything).Return(nil, nil)
	c.ReadDatabase("self_link")
	client.AssertCalled(t, "Read", "self_link", mock.Anything, mock.Anything)
}

func TestReadCollection(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	client.On("Read", "self_link", mock.Anything, mock.Anything).Return(nil, nil)
	c.ReadCollection("self_link")
	client.AssertCalled(t, "Read", "self_link", mock.Anything, mock.Anything)
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
	c := &DocumentDB{client, nil}
	client.On("Read", "self_link_doc", mock.Anything, mock.Anything).Return(nil, nil)
	c.ReadDocument("self_link_doc", &doc)
	client.AssertCalled(t, "Read", "self_link_doc", mock.Anything, mock.Anything)
}

func TestReadStoredProcedure(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	client.On("Read", "self_link", mock.Anything, mock.Anything).Return(nil, nil)
	c.ReadStoredProcedure("self_link")
	client.AssertCalled(t, "Read", "self_link", mock.Anything, mock.Anything)
}

func TestReadUserDefinedFunction(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	client.On("Read", "self_link", mock.Anything, mock.Anything).Return(nil, nil)
	c.ReadUserDefinedFunction("self_link")
	client.AssertCalled(t, "Read", "self_link", mock.Anything, mock.Anything)
}

func TestReadDatabases(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	client.On("Read", "dbs", mock.Anything, mock.Anything).Return(nil, nil)
	c.ReadDatabases()
	client.AssertCalled(t, "Read", "dbs", mock.Anything, mock.Anything)
}

func TestReadCollections(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	dbLink := "dblink/"
	client.On("Read", dbLink+"colls/", mock.Anything, mock.Anything).Return(nil, nil)
	c.ReadCollections(dbLink)
	client.AssertCalled(t, "Read", dbLink+"colls/", mock.Anything, mock.Anything)
}

func TestReadStoredProcedures(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	collLink := "colllink/"
	client.On("Read", collLink+"sprocs/", mock.Anything, mock.Anything).Return(nil, nil)
	c.ReadStoredProcedures(collLink)
	client.AssertCalled(t, "Read", collLink+"sprocs/", mock.Anything, mock.Anything)
}

func TestReadUserDefinedFunctions(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	collLink := "colllink/"
	client.On("Read", collLink+"udfs/", mock.Anything, mock.Anything).Return(nil, nil)
	c.ReadUserDefinedFunctions(collLink)
	client.AssertCalled(t, "Read", collLink+"udfs/", mock.Anything, mock.Anything)
}

func TestReadDocuments(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	collLink := "colllink/"
	client.On("Read", collLink+"docs/", mock.Anything, mock.Anything).Return(nil, nil)
	c.ReadDocuments(collLink, struct{}{})
	client.AssertCalled(t, "Read", collLink+"docs/", mock.Anything, mock.Anything)
}

func TestQueryDatabases(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	q := NewQuery("SELECT * FROM ROOT r")
	client.On("Query", "dbs", q).Return(nil)
	c.QueryDatabases(q)
	client.AssertCalled(t, "Query", "dbs", q)
}

func TestQueryCollections(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	q := NewQuery("SELECT * FROM ROOT r")
	client.On("Query", "db_self_link/colls/", q).Return(nil)
	c.QueryCollections("db_self_link/", q)
	client.AssertCalled(t, "Query", "db_self_link/colls/", q)
}

func TestQueryStoredProcedures(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	q := NewQuery("SELECT * FROM ROOT r")
	client.On("Query", "colls_self_link/sprocs/", q).Return(nil)
	c.QueryStoredProcedures("colls_self_link/", q)
	client.AssertCalled(t, "Query", "colls_self_link/sprocs/", q)
}

func TestQueryUserDefinedFunctions(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	q := NewQuery("SELECT * FROM ROOT r")
	client.On("Query", "colls_self_link/udfs/", q).Return(nil)
	c.QueryUserDefinedFunctions("colls_self_link/", q)
	client.AssertCalled(t, "Query", "colls_self_link/udfs/", q)
}

func TestQueryDocuments(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	collLink := "coll_self_link/"
	q := NewQuery("SELECT * FROM ROOT r")
	client.On("Query", collLink+"docs/", q).Return(nil)
	c.QueryDocuments(collLink, q, struct{}{})
	client.AssertCalled(t, "Query", collLink+"docs/", q)
}

func TestCreateDatabase(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	client.On("Create", "dbs", "{}").Return(nil)
	c.CreateDatabase("{}")
	client.AssertCalled(t, "Create", "dbs", "{}")
}

func TestCreateCollection(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	client.On("Create", "dbs/colls/", "{}").Return(nil)
	c.CreateCollection("dbs/", "{}")
	client.AssertCalled(t, "Create", "dbs/colls/", "{}")
}

func TestCreateStoredProcedure(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	client.On("Create", "dbs/colls/sprocs/", `{"id":"fn"}`).Return(nil)
	c.CreateStoredProcedure("dbs/colls/", `{"id":"fn"}`)
	client.AssertCalled(t, "Create", "dbs/colls/sprocs/", `{"id":"fn"}`)
}

func TestCreateUserDefinedFunction(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	client.On("Create", "dbs/colls/udfs/", `{"id":"fn"}`).Return(nil)
	c.CreateUserDefinedFunction("dbs/colls/", `{"id":"fn"}`)
	client.AssertCalled(t, "Create", "dbs/colls/udfs/", `{"id":"fn"}`)
}

func TestCreateDocument(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, defaultConfig}
	// TODO: test error situation, without id, etc...
	var doc Document
	client.On("Create", "dbs/colls/docs/", &doc).Return(nil)
	c.CreateDocument("dbs/colls/", &doc)
	client.AssertCalled(t, "Create", "dbs/colls/docs/", &doc)
	assert.NotEqual(t, doc.Id, "")
}

func TestUpsertDocument(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, defaultConfig}
	// TODO: test error situation, without id, etc...
	var doc Document
	client.On("Upsert", "dbs/colls/docs/", &doc).Return(nil)
	c.UpsertDocument("dbs/colls/", &doc)
	client.AssertCalled(t, "Upsert", "dbs/colls/docs/", &doc)
	assert.NotEqual(t, doc.Id, "")
}

func TestDeleteResource(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}

	client.On("Delete", "self_link_db").Return(nil)
	c.DeleteDatabase("self_link_db")
	client.AssertCalled(t, "Delete", "self_link_db")

	client.On("Delete", "self_link_coll").Return(nil)
	c.DeleteCollection("self_link_coll")
	client.AssertCalled(t, "Delete", "self_link_coll")

	client.On("Delete", "self_link_doc").Return(nil)
	c.DeleteDocument("self_link_doc")
	client.AssertCalled(t, "Delete", "self_link_doc")

	client.On("Delete", "self_link_sproc").Return(nil)
	c.DeleteDocument("self_link_sproc")
	client.AssertCalled(t, "Delete", "self_link_sproc")

	client.On("Delete", "self_link_udf").Return(nil)
	c.DeleteDocument("self_link_udf")
	client.AssertCalled(t, "Delete", "self_link_udf")
}

func TestReplaceDatabase(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	client.On("Replace", "db_link", "{}").Return(nil)
	c.ReplaceDatabase("db_link", "{}")
	client.AssertCalled(t, "Replace", "db_link", "{}")
}

func TestReplaceDocument(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	client.On("Replace", "doc_link", "{}").Return(nil)
	c.ReplaceDocument("doc_link", "{}")
	client.AssertCalled(t, "Replace", "doc_link", "{}")
}

func TestReplaceStoredProcedure(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	client.On("Replace", "sproc_link", "{}").Return(nil)
	c.ReplaceStoredProcedure("sproc_link", "{}")
	client.AssertCalled(t, "Replace", "sproc_link", "{}")
}

func TestReplaceUserDefinedFunction(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	client.On("Replace", "udf_link", "{}").Return(nil)
	c.ReplaceUserDefinedFunction("udf_link", "{}")
	client.AssertCalled(t, "Replace", "udf_link", "{}")
}

func TestExecuteStoredProcedure(t *testing.T) {
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	client.On("Execute", "sproc_link", "{}").Return(nil)
	c.ExecuteStoredProcedure("sproc_link", "{}", struct{}{})
	client.AssertCalled(t, "Execute", "sproc_link", "{}")
}

func TestQueryPartitionKeyRanges(t *testing.T) {
	expectedRanges := []PartitionKeyRange{
		PartitionKeyRange{
			PartitionKeyRangeID: "1",
		},
	}
	client := &ClientStub{}
	c := &DocumentDB{client, nil}
	client.On("Read", "coll_link/pkranges/", mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		r := args.Get(1).(*queryPartitionKeyRangesRequest)
		r.Ranges = expectedRanges
	}).Return(&Response{}, nil)
	ranges, err := c.QueryPartitionKeyRanges("coll_link/", nil)
	client.AssertCalled(t, "Read", "coll_link/pkranges/", mock.Anything, mock.Anything)
	assert.NoError(t, err)
	assert.Equal(t, expectedRanges, ranges, "Ranges are different")
}
