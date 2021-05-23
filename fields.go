package structfields

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
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
		ss = append(ss, Structs(pkg)...)
	}
	return ss, nil
}

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

func ResolveType(pkg *packages.Package, name string) (*ast.StructType, *ast.File, bool) {
	for _, file := range pkg.Syntax {
		var stype *ast.StructType
		ast.Inspect(file, func(n ast.Node) bool {
			if stype != nil {
				return false
			}
			spec, ok := n.(*ast.TypeSpec)
			if !ok {
				return true
			}
			if spec.Name.Name != name {
				return false
			}
			st, ok := spec.Type.(*ast.StructType)
			if ok {
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

func Fields(pkg *packages.Package, file *ast.File, stype *ast.StructType) []*FieldType {
	var ff []*FieldType
	for _, f := range stype.Fields.List {
		// if there are no names, it's embedded
		if len(f.Names) == 0 {
			selexpr, ok := f.Type.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			pkg0, ok := ResolvePackage(pkg, file, exprfmt(selexpr.X))
			if !ok {
				continue
			}
			if stype0, file, ok := ResolveType(pkg0, exprfmt(selexpr.Sel)); ok {
				ff = append(ff, Fields(pkg0, file, stype0)...)
			}
			continue
		}
		for _, name := range f.Names {
			if !ast.IsExported(name.Name) {
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
				Name: name.Name,
				Type: exprfmt(f.Type),
				Doc:  doctxt,
				Tag:  tagtxt,
			})
		}
	}
	return ff
}

func Structs(pkg *packages.Package) []*StructType {
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
			s.Fields = Fields(pkg, file, stype)
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
