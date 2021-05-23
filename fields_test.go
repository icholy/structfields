package structfields

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestLoad(t *testing.T) {
	ss, err := Load("testmod", ".")
	assert.NilError(t, err)
	lookup := func(name string) *StructType {
		t.Helper()
		for _, s := range ss {
			if s.Name == name {
				return s
			}
		}
		t.Fatalf("struct not found: %q", name)
		return nil // unreachable
	}
	t.Run("A", func(t *testing.T) {
		assert.DeepEqual(t, lookup("A"), &StructType{
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
		t.SkipNow()
		assert.DeepEqual(t, lookup("B"), &StructType{
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
