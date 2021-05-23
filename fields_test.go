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

func TestLoad(t *testing.T) {
	ss := Find(load(t))
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

func TestResolve(t *testing.T) {
	pkg := load(t)
	t.Run("A", func(t *testing.T) {
		_, ok := Resolve(pkg, "A")
		assert.Assert(t, ok)
	})
	t.Run("sub.C", func(t *testing.T) {
		_, ok := Resolve(pkg, "sub.C")
		assert.Assert(t, ok)
	})
}

func TestFields(t *testing.T) {
	pkg := load(t)
	stype, ok := Resolve(pkg, "B")
	assert.Assert(t, ok)
	ff := Fields(pkg, stype)
	assert.DeepEqual(t, ff, []*FieldType{
		{
			Name: "F1",
			Type: "string",
		},
	})
}
