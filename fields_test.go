package structfields

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestLoad(t *testing.T) {
	ss, err := Load("testmod", "testmod")
	assert.NilError(t, err)
	t.Run("A", func(t *testing.T) {
		s := lookup(t, ss, "A")
		assert.DeepEqual(t, s, &StructType{
			Name: "A",
			Fields: []*FieldType{
				{Name: "F1", Type: "string"},
				{Name: "F2", Type: "int"},
				{Name: "F3", Type: "bool"},
				{Name: "F4", Type: "time.Duration"},
			},
		})
	})
	t.Run("B", func(t *testing.T) {
		s := lookup(t, ss, "B")
		assert.DeepEqual(t, s, &StructType{
			Name: "B",
			Fields: []*FieldType{
				{Name: "F1", Type: "string"},
				{Name: "F2", Type: "int"},
				{Name: "F3", Type: "bool"},
				{Name: "F4", Type: "time.Duration"},
			},
		})
	})
}

func lookup(t *testing.T, ss []*StructType, name string) *StructType {
	t.Helper()
	for _, s := range ss {
		if s.Name == name {
			return s
		}
	}
	t.Fatalf("struct not found: %q", name)
	return nil // unreachable
}
