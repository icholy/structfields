package structfields

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/packages"
)

type StructType struct {
	Doc    string
	Group  string
	Name   string
	Fields []*FieldType
}

type FieldType struct {
	Name string
	Type string
	Doc  string
	Tag  string
}

func Load(dir string, pkgpath string) ([]*StructType, error) {
	cfg := &packages.Config{
		Dir: dir,
		Mode: packages.NeedName |
			packages.NeedDeps |
			packages.NeedFiles |
			packages.NeedSyntax |
			packages.NeedImports,
	}
	pkgs, err := packages.Load(cfg, pkgpath)
	if err != nil {
		return nil, err
	}
	var ss []*StructType
	for _, pkg := range pkgs {
		ss = append(ss, Find(pkg)...)
	}
	return ss, nil
}

func ResolveType(pkg *packages.Package, name string) (*ast.TypeSpec, bool) {
	// if the name is prefixed by a path, change pkg to the appropriate one
	if i := strings.IndexByte(name, '.'); i >= 0 {
		pkgname := name[:i]
		name = name[i+1:]
		var found bool
		for _, pkg0 := range pkg.Imports {
			if pkg0.Name == pkgname {
				pkg = pkg0
				found = true
				break
			}
		}
		if !found {
			return nil, false
		}
	}
	for _, file := range pkg.Syntax {
		var spec *ast.TypeSpec
		ast.Inspect(file, func(n ast.Node) bool {
			if spec != nil {
				return false
			}
			s, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}
			if s.Name.Name == name {
				spec = s
				return false
			}
			return true
		})
		if spec != nil {
			return spec, true
		}
	}
	return nil, false
}

func StructFields(pkg *packages.Package, stype *ast.StructType) []*FieldType {
	var ff []*FieldType
	for _, f := range stype.Fields.List {
		if len(f.Names) != 1 {
			continue
		}
		name := f.Names[0].Name
		if !ast.IsExported(name) {
			continue
		}
		var doctxt string
		if f.Doc != nil {
			doctxt = f.Doc.Text()
		}
		if doctxt == "" && f.Comment != nil {
			doctxt = f.Comment.Text()
		}
		var tagtxt string
		if f.Tag != nil {
			tagtxt = f.Tag.Value
		}
		ff = append(ff, &FieldType{
			Name: name,
			Type: exprfmt(f.Type),
			Doc:  doctxt,
			Tag:  tagtxt,
		})
	}
	return ff
}

func Find(pkg *packages.Package) []*StructType {
	var ss []*StructType
	for _, file := range pkg.Syntax {
		ast.Inspect(file, func(n ast.Node) bool {
			var s StructType
			decl, ok := n.(*ast.GenDecl)
			if !ok || decl.Tok != token.TYPE || len(decl.Specs) != 1 {
				return true
			}
			spec, ok := decl.Specs[0].(*ast.TypeSpec)
			if !ok {
				return true
			}
			s.Name = spec.Name.String()
			stype, ok := spec.Type.(*ast.StructType)
			if !ok {
				return true
			}
			if decl.Doc != nil {
				s.Doc = decl.Doc.Text()
				for _, comment := range decl.Doc.List {
					text := comment.Text
					i := strings.Index(text, "go:docgen")
					if i >= 0 {
						s.Group = strings.TrimSpace(text[i+9:])
						break
					}
				}
			}
			s.Fields = StructFields(pkg, stype)
			ss = append(ss, &s)
			return false
		})
	}
	return ss
}

func exprfmt(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.StarExpr:
		return "*" + exprfmt(e.X)
	case *ast.Ident:
		return e.Name
	case *ast.ArrayType:
		return "[]" + exprfmt(e.Elt)
	case *ast.MapType:
		return fmt.Sprintf("map[%s]%s", exprfmt(e.Key), exprfmt(e.Value))
	case *ast.SelectorExpr:
		return fmt.Sprintf("%s.%s", exprfmt(e.X), exprfmt(e.Sel))
	case *ast.StructType:
		return "struct{ ... }"
	case *ast.FuncType:
		return "func(...) ..."
	default:
		panic(fmt.Errorf("not implemented: %#v", expr))
	}
}
