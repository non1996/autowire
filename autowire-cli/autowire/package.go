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
	"path"
	"strings"

	"github.com/non1996/go-jsonobj/container"

	"github.com/non1996/go-autowire/autowire-cli/annotation"
	"github.com/non1996/go-autowire/autowire-cli/internal/assert"
)

type Package struct {
	Config
	ast          *ast.Package
	AbsolutePath string // 包绝对路径
	RelativePath string // 包在项目下的相对路径
	Package      string
	Imports      Imports
	genFile      []byte
	Components   []*Component
	Applications []*Application
}

func NewPackage(
	config Config,
	dir string,
	ast *ast.Package,
) *Package {
	return &Package{
		Config:       config,
		ast:          ast,
		AbsolutePath: dir,
		RelativePath: strings.ReplaceAll(dir, config.Root, ""),
		Package:      ast.Name,
	}
}

func parsePackage(
	conf Config,
	dir string,
) *Package {
	fSet := token.NewFileSet()
	packages, err := parser.ParseDir(fSet,
		dir,
		func(info fs.FileInfo) bool {
			return info.Name() == conf.AutowireFileName
		},
		parser.DeclarationErrors,
	)
	if err != nil {
		panic(err)
	}
	if len(packages) == 0 {
		return nil
	}

	assert.Assert(len(packages) == 1)

	pkg := NewPackage(conf, dir, packages[container.MapKeys(packages)[0]])

	file := pkg.ast.Files[path.Join(dir, conf.AutowireFileName)]

	pkg.Imports = parseImport(file.Decls[0].(*ast.GenDecl))

	var annos []annotation.PrimaryAnnotation
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		annos = append(annos, annotation.Parse("a", genDecl)...)
	}

	for _, anno := range annos {
		if anno.Name == "Component" {
			component := parseComponent(pkg.Module, pkg.RelativePath, anno)
			pkg.Components = append(pkg.Components, component)
			fmt.Printf("component %+v\n", component)
		} else if anno.Name == "Application" {
			app := parseApplication(anno)
			pkg.Applications = append(pkg.Applications, app)
			fmt.Printf("app %+v\n", app)
		}
	}

	pkg.Imports.Add(Import{
		Alias:       "autowire",
		PackagePath: "github.com/non1996/go-autowire/autowire",
	})
	pkg.Imports.Map.Remove("github.com/non1996/go-autowire/a")

	for _, app := range pkg.Applications {
		for _, configuration := range app.Configurations {
			pkg.Imports.Add(Import{
				Alias:       "_",
				PackagePath: configuration,
				HasAlias:    true,
			})
		}
	}

	pkg.GenFileName = conf.GenFileName

	return pkg
}

func (p *Package) generateFactories() {
	if len(p.Components) == 0 && len(p.Applications) == 0 {
		return
	}
	var buffer = bytes.NewBuffer(nil)

	err := factoryTemplate.ExecuteTemplate(buffer, "tmplGenFile", p)
	if err != nil {
		panic(err)
	}
	p.genFile = buffer.Bytes()
	p.genFile, err = format.Source(buffer.Bytes())
	if err != nil {
		panic(err)
	}
}

func (p *Package) output() {
	if len(p.Components) == 0 && len(p.Applications) == 0 {
		return
	}
	err := os.WriteFile(path.Join(p.AbsolutePath, p.Config.GenFileName), p.genFile, os.ModePerm)
	if err != nil {
		panic(err)
	}
}
