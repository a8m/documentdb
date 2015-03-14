package documentdb

import (
	"testing"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/assert"
)

type ClientStub struct {
	mock.Mock
}

func (c *ClientStub) Read(link string, ret interface{}) error {
	c.Called(link)
	return nil
}

func TestNew(t *testing.T) {
	assert := assert.New(t)
	client := New("url", Config{"config"})
	assert.IsType(client, &DocumentDB{}, "Should return DocumentDB object")
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
