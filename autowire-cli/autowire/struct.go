package autowire

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/non1996/go-jsonobj/container"
)

type Struct struct {
	Package *Package
	Name    string
	Fields  *container.OrderedMap[string, *Field]
	Methods *container.OrderedMap[string, *Method]
	astSpec *ast.TypeSpec
}

func NewStruct(p *Package, name string) *Struct {
	return &Struct{
		Package: p,
		Name:    name,
		Fields:  container.NewOrderedMap[string, *Field](),
		Methods: container.NewOrderedMap[string, *Method](),
		astSpec: nil,
	}
}

func (s *Struct) GetField(name string) *Field {
	return s.Fields.Get(name)
}

func (s *Struct) GetMethod(name string) *Method {
	return s.Methods.Get(name)
}

type Field struct {
	Name      string
	Type      Type
	Anonymous bool
}

type Method struct {
	Name   string
	Param  *container.OrderedMap[string, *Field]
	Result *container.OrderedMap[string, *Field]
	ast    *ast.FuncDecl
}

func (m *Method) GetFirstParam() (field *Field) {
	m.Param.Foreach(func(s string, f *Field) {
		if field == nil {
			field = f
		}
	})

	return
}

func (m *Method) GetFirstResult() (field *Field) {
	m.Result.Foreach(func(s string, f *Field) {
		if field == nil {
			field = f
		}
	})

	return
}

func (p *Package) getStruct(name string) *Struct {
	strct := p.Structs.Get(name)
	if strct == nil {
		strct = NewStruct(p, name)
		p.Structs.Add(strct.Name, strct)
	}

	return strct
}

func (p *Package) getStructDecl(name string) *Struct {
	return p.Structs.Get(name)
}

func (p *Package) parseStructs() {
	for path, file := range p.ast.Files {
		if strings.HasSuffix(path, p.AutowireFileName) {
			continue
		}

		imports := parseImport(file.Imports)

		p.Imports.Merge(imports)

		for _, decl := range file.Decls {
			switch t := decl.(type) {
			case *ast.GenDecl:
				if t.Tok != token.TYPE {
					continue
				}

				for _, spec := range t.Specs {
					typeSpec := spec.(*ast.TypeSpec)
					structType, ok := spec.(*ast.TypeSpec).Type.(*ast.StructType)
					if !ok {
						continue
					}

					strct := p.getStruct(typeSpec.Name.Name)
					strct.astSpec = typeSpec
					strct.Fields = p.convertFieldList(imports, structType.Fields)
				}
			case *ast.FuncDecl:
				if t.Recv == nil {
					continue
				}

				typ := p.NewType(t.Recv.List[0].Type)
				strct := p.getStruct(typ.TypeName())

				method := &Method{
					Name: t.Name.Name,
					ast:  t,
				}

				method.Param = p.convertFieldList(imports, t.Type.Params)
				method.Result = p.convertFieldList(imports, t.Type.Results)

				strct.Methods.Add(method.Name, method)
			default:
				continue
			}
		}
	}
}

func (p *Package) convertFieldList(
	imports *Imports,
	fields *ast.FieldList,
) (list *container.OrderedMap[string, *Field]) {
	list = container.NewOrderedMap[string, *Field]()

	if fields == nil {
		return list
	}

	for idx, field := range fields.List {
		typ := p.NewType(field.Type)
		if !typ.isThisPackage() {
			typ.SetImport(imports.GetByAlias(typ.getPackage()).PackagePath)
		}

		if len(field.Names) == 0 {
			list.Add(fmt.Sprintf("_%d", idx), &Field{
				Type:      typ,
				Anonymous: true,
			})
			continue
		}

		for _, name := range field.Names {
			list.Add(name.Name, &Field{
				Name: name.Name,
				Type: typ,
			})
		}
	}

	return list
}
