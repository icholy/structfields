package structfields

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"
)

// StructType contains details about a struct and it's fields.
type StructType struct {
	Doc        string
	Name       string
	Directives []string
	Fields     []*FieldType
}

// FieldType contains information about a struct field.
type FieldType struct {
	Name    string
	Type    string
	Doc     string
	Comment string
	Tag     string
}

// Needs is the packages.Mode with all the flags required to find fields.
const Needs = packages.NeedName |
	packages.NeedDeps |
	packages.NeedFiles |
	packages.NeedSyntax |
	packages.NeedImports

// Load finds all structs in the packages specified in the patterns.
// Embeded fields are treated the same as regular fields.
func Load(dir string, patterns ...string) ([]*StructType, error) {
	cfg := &packages.Config{
		Dir:  dir,
		Mode: Needs,
	}
	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, err
	}
	var ss []*StructType
	for _, pkg := range pkgs {
		ss = append(ss, Structs(pkg)...)
	}
	return ss, nil
}

// ResolvePackage will resolve a package by its name.
// It's legal to pass a nil file, but then import aliases and duplicates are not handled.
func ResolvePackage(pkg *packages.Package, file *ast.File, name string) (*packages.Package, bool) {
	if file == nil {
		// if there's no file, assume there are no aliases
		for _, pkg0 := range pkg.Imports {
			if pkg0.Name == name {
				return pkg0, true
			}
		}
		return nil, false
	}
	for _, imp := range file.Imports {
		pkgpath, _ := strconv.Unquote(imp.Path.Value)
		if imp.Name != nil && imp.Name.Name == name {
			pkg0, ok := pkg.Imports[pkgpath]
			return pkg0, ok
		}
		if pkg0, ok := pkg.Imports[pkgpath]; ok && pkg0.Name == name {
			return pkg0, true
		}
	}
	return nil, false
}

// ResolveType returns a struct type defined in the provided package along with the file it's declared in.
// The bool return value will be false if the name could not be resolved.
func ResolveType(pkg *packages.Package, name string) (*ast.StructType, *ast.File, bool) {
	for _, file := range pkg.Syntax {
		var stype *ast.StructType
		ast.Inspect(file, func(n ast.Node) bool {
			if stype != nil {
				return false
			}
			tspec, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}
			if tspec.Name.Name != name {
				return false
			}
			if st, ok := tspec.Type.(*ast.StructType); ok {
				stype = st
				return false
			}
			return true
		})
		if stype != nil {
			return stype, file, true
		}
	}
	return nil, nil, false
}

// Fields returns a list of the struct type's fields.
// A nil file may be passed, but this limits the ability to resolve embeded types.
func Fields(pkg *packages.Package, file *ast.File, stype *ast.StructType) []*FieldType {
	ff := []*FieldType{}
	for _, f := range stype.Fields.List {
		// if there are no names, it's embedded
		if len(f.Names) == 0 {
			pkg0 := pkg
			expr := f.Type
			// if it's a selector expression, it's a type in a different package
			if selexpr, ok := expr.(*ast.SelectorExpr); ok {
				if pkg0, ok = ResolvePackage(pkg, file, exprfmt(selexpr.X)); !ok {
					continue
				}
				expr = selexpr.Sel
			}
			if stype0, file, ok := ResolveType(pkg0, exprfmt(expr)); ok {
				ff = append(ff, Fields(pkg0, file, stype0)...)
			}
			continue
		}
		for _, name := range f.Names {
			if !ast.IsExported(name.Name) {
				continue
			}
			ft := FieldType{
				Name: name.String(),
				Type: exprfmt(f.Type),
			}
			if f.Doc != nil {
				ft.Doc = f.Doc.Text()
			}
			if f.Comment != nil {
				ft.Comment = f.Comment.Text()
			}
			if f.Tag != nil {
				ft.Tag = f.Tag.Value
			}
			ff = append(ff, &ft)
		}
	}
	return ff
}

// Structs finds all structs in the provided package.
// Embeded fields are treated the same as regular fields.
func Structs(pkg *packages.Package) []*StructType {
	ss := []*StructType{}
	for _, file := range pkg.Syntax {
		ast.Inspect(file, func(n ast.Node) bool {
			var s StructType
			decl, ok := n.(*ast.GenDecl)
			if !ok || decl.Tok != token.TYPE {
				return true
			}
			for _, spec := range decl.Specs {
				tspec, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}
				s.Name = tspec.Name.String()
				stype, ok := tspec.Type.(*ast.StructType)
				if !ok {
					continue
				}
				s.Directives = []string{}
				if decl.Doc != nil {
					s.Doc = decl.Doc.Text()
					for _, comment := range decl.Doc.List {
						if strings.HasPrefix(comment.Text, "//go:") {
							s.Directives = append(s.Directives, comment.Text[5:])
						}
					}
				}
				s.Fields = Fields(pkg, file, stype)
				ss = append(ss, &s)
			}
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
	case *ast.ChanType:
		return fmt.Sprintf("chan %s", exprfmt(e.Value))
	default:
		panic(fmt.Errorf("not implemented: %#v", expr))
	}
}
