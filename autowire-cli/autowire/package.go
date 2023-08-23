package autowire

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/non1996/go-jsonobj/container"

	"github.com/non1996/go-autowire/autowire-cli/internal/assert"
)

type Package struct {
	Config
	fset         *token.FileSet
	ast          *ast.Package
	AbsolutePath string // 包绝对路径
	PackagePath  string // 包在项目下的相对路径
	PackageName  string // 包名
	Imports      *Imports
	Structs      *container.OrderedMap[string, *Struct]
	AnnoFile     *AnnoFile // 注解文件
	GenFile      *GenFile
}

func NewPackage(
	config Config,
	dir string,
	fset *token.FileSet,
	ast *ast.Package,
) *Package {
	return &Package{
		Config:       config,
		fset:         fset,
		ast:          ast,
		AbsolutePath: dir,
		PackagePath:  strings.ReplaceAll(dir, config.Root, ""),
		PackageName:  ast.Name,
		Imports:      NewImports(),
		Structs:      container.NewOrderedMap[string, *Struct](),
	}
}

func parsePackage(
	conf Config,
	dir string,
) *Package {
	fset := token.NewFileSet()
	astPkgs, err := parser.ParseDir(
		fset,
		dir,
		func(info fs.FileInfo) bool {
			return info.Name() != conf.GenFileName
		},
		parser.DeclarationErrors,
	)
	if err != nil {
		panic(err)
	}
	if len(astPkgs) == 0 {
		return nil
	}
	assert.Assert(len(astPkgs) == 1, "more than one package in a directory")

	astPkg := astPkgs[container.MapKeys(astPkgs)[0]]

	pkg := NewPackage(conf, dir, fset, astPkg)

	pkg.parseAnnoFile()

	pkg.parseStructs()

	return pkg
}

func (p *Package) evaluate() {
	if p.AnnoFile == nil {
		return
	}

	p.GenFile = &GenFile{
		Package:      p,
		AbsolutePath: filepath.Join(p.AbsolutePath, p.Config.GenFileName),
		Imports:      NewImports(),
	}

	p.GenFile.Imports.Add(&Import{
		Alias:       "autowire",
		PackagePath: "github.com/non1996/go-autowire/autowire",
	})

	for _, anno := range p.AnnoFile.Annotations {
		if anno.Name == "Component" {
			component := p.evaluateComponent(anno)
			p.GenFile.Components = append(p.GenFile.Components, component)
			fmt.Printf("component %+v\n", component)
		} else if anno.Name == "Application" {
			app := p.evaluateApplication(anno)
			p.GenFile.Applications = append(p.GenFile.Applications, app)
			fmt.Printf("app %+v\n", app)
		}
	}

	for _, component := range p.GenFile.Components {
		p.link(component)
	}

	for _, application := range p.GenFile.Applications {
		p.link(&application.Component)
	}
}

func (p *Package) link(component *Component) {
	strct := p.getStructDecl(component.Type.TypeName())
	assert.Assert(strct != nil)

	for _, injector := range component.Injectors {
		baseInjector := injector.Base()
		baseInjector.IsMethod = strct.GetField(baseInjector.Value) == nil

		var fieldType Type

		if baseInjector.IsMethod {
			method := strct.GetMethod(baseInjector.Value)
			assert.Assert(method.Param.Size() == 1)
			fieldType = method.GetFirstParam().Type
		} else {
			field := strct.GetField(baseInjector.Value)
			fieldType = field.Type
		}

		imp := p.Imports.GetByPath(fieldType.i)
		if imp != nil {
			p.GenFile.Imports.Add(imp)
		}

		if injector.Kind() == InjectorKindComponent {
			if fieldType.isSlice() {
				fieldType = fieldType.sliceElem()
				injector.(*ComponentInjector).IsSlice = true
			}

			injector.(*ComponentInjector).Type = fieldType
		} else if injector.Kind() == InjectorKindEnv {
			assert.Assert(fieldType.TypeName() == "string")
		}
	}

	for _, bean := range component.Beans {
		method := strct.GetMethod(bean.Method)
		assert.Assert(method.Result.Size() == 1)

		fieldType := method.GetFirstResult().Type

		imp := p.Imports.GetByPath(fieldType.i)
		if imp != nil {
			p.GenFile.Imports.Add(imp)
		}

		bean.Type = fieldType
		bean.Ptr = fieldType.isPointer()
	}

	for _, property := range component.Properties {
		property.IsMethod = strct.GetField(property.Value) == nil

		var field *Field

		if property.IsMethod {
			method := strct.GetMethod(property.Value)
			assert.Assert(method.Result.Size() == 0)
			assert.Assert(method.Result.Size() == 1)
			field = method.GetFirstResult()
			property.Type = field.Type
		} else {
			field = strct.GetField(property.Value)
			property.Type = field.Type
		}

		imp := p.Imports.GetByPath(field.Type.i)
		if imp != nil {
			p.GenFile.Imports.Add(imp)
		}
	}

	if component.PostConstruct != nil {
		method := strct.GetMethod(component.PostConstruct.MethodName)
		component.PostConstruct.HasErrorResp = method.Result.Size() == 1 &&
			method.GetFirstResult().Type.TypeName() == "error"
	}
}

func (p *Package) format() {
	if p.GenFile == nil {
		return
	}
	var buffer = bytes.NewBuffer(nil)

	err := factoryTemplate.ExecuteTemplate(buffer, "tmplGenFile", p.GenFile)
	if err != nil {
		panic(err)
	}
	p.GenFile.Output = buffer.Bytes()
	p.GenFile.Output, err = format.Source(buffer.Bytes())
	if err != nil {
		panic(err)
	}
}

func (p *Package) output() {
	if p.GenFile == nil {
		return
	}
	err := os.WriteFile(p.GenFile.AbsolutePath, p.GenFile.Output, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
