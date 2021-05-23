package structfields

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestLoad(t *testing.T) {
	ss, err := Load("testmod", "testmod")
	assert.NilError(t, err)
	assert.DeepEqual(t, ss, []*StructType{
		{
			Name: "A",
			Fields: []*FieldType{
				{Name: "F1", Type: "string"},
				{Name: "F2", Type: "int"},
				{Name: "F3", Type: "bool"},
			},
		},
	})
}
