package structfields

import (
	"testing"

	"golang.org/x/tools/go/packages"
	"gotest.tools/v3/assert"
)

func load(t *testing.T) *packages.Package {
	t.Helper()
	cfg := &packages.Config{
		Dir: "testmod",
		Mode: packages.NeedName |
			packages.NeedDeps |
			packages.NeedFiles |
			packages.NeedSyntax |
			packages.NeedImports,
	}
	pkgs, err := packages.Load(cfg, ".")
	assert.NilError(t, err)
	assert.Equal(t, len(pkgs), 1, "expecting one package")
	return pkgs[0]
}

func TestStructs(t *testing.T) {
	ss := Structs(load(t))
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
				{Name: "F5", Type: "byte"},
				{Name: "F6", Type: "byte"},
			},
		})
	})
	t.Run("B", func(t *testing.T) {
		assert.DeepEqual(t, lookup("B"), &StructType{
			Name: "B",
			Fields: []*FieldType{
				{Name: "F1", Type: "string"},
				{Name: "F42", Type: "int64"},
			},
		})
	})
	t.Run("E", func(t *testing.T) {
		assert.DeepEqual(t, lookup("E"), &StructType{
			Name: "E",
			Doc:  "E is a struct\n",
			Fields: []*FieldType{
				{
					Name:    "F1",
					Type:    "string",
					Doc:     "F1 is a string\n",
					Comment: "F1 is the first field\n",
					Tag:     "`json:\"Test\"`",
				},
				{
					Name:    "F2",
					Type:    "int",
					Comment: "F2 only has a comment\n",
				},
				{
					Name: "F3",
					Type: "bool",
					Doc:  "F3 only has a doc\n",
				},
			},
		})
	})
}

func TestResolve(t *testing.T) {
	pkg := load(t)
	t.Run("A", func(t *testing.T) {
		_, _, ok := ResolveType(pkg, "A")
		assert.Assert(t, ok)
	})
	t.Run("sub.C", func(t *testing.T) {
		pkg0, ok := ResolvePackage(pkg, nil, "sub")
		assert.Assert(t, ok)
		_, _, ok = ResolveType(pkg0, "C")
		assert.Assert(t, ok)
	})
}

func TestFields(t *testing.T) {
	pkg := load(t)
	t.Run("D", func(t *testing.T) {
		stype, file, ok := ResolveType(pkg, "D")
		assert.Assert(t, ok)
		ff := Fields(pkg, file, stype)
		assert.DeepEqual(t, ff, []*FieldType{
			{
				Name: "F1",
				Type: "string",
			},
		})
	})
}
