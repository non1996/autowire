package autowire

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/non1996/go-jsonobj/container"

	"github.com/non1996/go-autowire/autowire-cli/internal/assert"
)

// Imports pkg -> import
type Imports struct {
	Map      *container.OrderedMap[string, Import]
	AliasMap map[string]Import
}

func (i *Imports) List() (res []Import) {
	i.Map.Foreach(func(s string, i Import) {
		res = append(res, i)
	})

	return res
}

func (i *Imports) Add(imp Import) {
	i.Map.Add(imp.PackagePath, imp)
	if imp.Alias != "_" && imp.Alias != "" && imp.Alias != "." {
		i.AliasMap[imp.Alias] = imp
	}
}

func (i *Imports) RemoveByPath(path string) {
	if !i.Map.Exist(path) {
		return
	}
	imp := i.Map.Get(path)
	i.Map.Remove(path)
	if imp.Alias != "_" && imp.Alias != "" && imp.Alias != "." {
		delete(i.AliasMap, imp.Alias)
	}
}

type Import struct {
	Alias       string
	PackagePath string
	HasAlias    bool
}

func parseImport(decl *ast.GenDecl) Imports {
	assert.Assert(decl.Tok == token.IMPORT)

	imports := Imports{
		Map:      container.NewOrderedMap[string, Import](),
		AliasMap: map[string]Import{},
	}

	for _, spec := range decl.Specs {
		importSpec := spec.(*ast.ImportSpec)

		i := Import{}

		i.PackagePath = mustStringLit(importSpec.Path)

		if importSpec.Name != nil {
			i.Alias = importSpec.Name.Name
			i.HasAlias = true
		} else {
			i.Alias = i.PackagePath[strings.LastIndex(i.PackagePath, "/")+1:]
		}

		imports.Add(i)
	}

	return imports
}
