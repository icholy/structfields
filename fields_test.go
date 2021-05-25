package structfields

import (
	"testing"

	"github.com/google/go-cmp/cmp/cmpopts"
	"golang.org/x/tools/go/packages"
	"gotest.tools/v3/assert"
)

var ignoreType = cmpopts.IgnoreFields(FieldType{}, "Type")

func load(t *testing.T) *packages.Package {
	t.Helper()
	cfg := &packages.Config{
		Dir:  "testmod",
		Mode: Needs,
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
			Name:       "A",
			Directives: []string{},
			Fields: []*FieldType{
				{Name: "F1"},
				{Name: "F2"},
				{Name: "F3"},
				{Name: "F4"},
				{Name: "F5"},
				{Name: "F6"},
			},
		}, ignoreType)
	})
	t.Run("B", func(t *testing.T) {
		assert.DeepEqual(t, lookup("B"), &StructType{
			Name:       "B",
			Directives: []string{},
			Fields: []*FieldType{
				{Name: "F1"},
				{Name: "F42"},
			},
		}, ignoreType)
	})
	t.Run("E", func(t *testing.T) {
		assert.DeepEqual(t, lookup("E"), &StructType{
			Name:       "E",
			Doc:        "E is a struct\n",
			Directives: []string{"what", "directive"},
			Fields: []*FieldType{
				{
					Name:    "F1",
					Doc:     "F1 is a string\n",
					Comment: "F1 is the first field\n",
					Tag:     "`json:\"Test\"`",
				},
				{
					Name:    "F2",
					Comment: "F2 only has a comment\n",
				},
				{
					Name: "F3",
					Doc:  "F3 only has a doc\n",
				},
			},
		}, ignoreType)
	})
	t.Run("F", func(t *testing.T) {
		assert.DeepEqual(t, lookup("F"), &StructType{
			Name:       "F",
			Directives: []string{},
			Fields: []*FieldType{
				{Name: "F1"},
				{Name: "F2"},
				{Name: "F3"},
				{Name: "F4"},
				{Name: "F5"},
				{Name: "F6"},
			},
		}, ignoreType)
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
			},
		}, ignoreType)
	})
}
