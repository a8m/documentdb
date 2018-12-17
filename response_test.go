package documentdb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpectStatusCode(t *testing.T) {

	expecations := []struct {
		status  int
		result  bool
		message string
	}{
		{200, true, "tesing 200, should be true"},
		{400, false, "tesing 400, should be false"},
	}

	for _, e := range expecations {
		actual := expectStatusCode(200)(e.status)
		assert.Equal(t, e.result, actual, e.message)
	}

}

func TestExpectStatusCodeXX(t *testing.T) {

	expecations := []struct {
		status  int
		result  bool
		message string
	}{
		{199, false, "bellow range"},
		{200, true, "range begining"},
		{250, true, "in range"},
		{299, true, "range end"},
		{300, false, "above range"},
	}

	for _, e := range expecations {
		actual := expectStatusCodeXX(200)(e.status)
		assert.Equal(t, e.result, actual, e.message)
	}

}
