package structfields

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestLoad(t *testing.T) {
	ss, err := Load("testdata", "testdata")
	assert.NilError(t, err)
	assert.DeepEqual(t, ss, []*StructType{
		{},
	})
}
