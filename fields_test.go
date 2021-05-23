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
		spec, ok := ResolveType(pkg, "A")
		assert.Assert(t, ok)
		assert.Equal(t, spec.Name.String(), "A")
	})
}
