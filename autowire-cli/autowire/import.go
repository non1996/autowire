package autowire

import (
	"go/ast"
	"strings"

	"github.com/non1996/go-jsonobj/container"
)

type Imports struct {
	Map      *container.OrderedMap[string, *Import] // pkg -> import
	AliasMap map[string]string                      // alias -> pkg
}

func NewImports() *Imports {
	return &Imports{
		Map:      container.NewOrderedMap[string, *Import](),
		AliasMap: map[string]string{},
	}
}

func (i *Imports) List() (res []*Import) {
	i.Map.Foreach(func(s string, i *Import) {
		res = append(res, i)
	})

	return res
}

func (i *Imports) Add(imp *Import) {
	i.Map.Add(imp.PackagePath, imp)
	if imp.Alias != "_" && imp.Alias != "" && imp.Alias != "." {
		i.AliasMap[imp.Alias] = imp.PackagePath
	}
}

func (i *Imports) AddIfAbsent(imp *Import) {
	if i.Map.Exist(imp.PackagePath) {
		return
	}

	i.Add(imp)
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

func (i *Imports) GetByAlias(alias string) *Import {
	path, exist := i.AliasMap[alias]
	if !exist {
		return nil
	}

	return i.Map.Get(path)
}

func (i *Imports) GetByPath(path string) *Import {
	return i.Map.Get(path)
}

func (i *Imports) Merge(imports *Imports) {
	for _, imp := range imports.List() {
		if i.Map.Exist(imp.PackagePath) {
			continue
		}

		alias := imp.Alias
		for container.MapExist(i.AliasMap, alias) {
			alias = alias + "_"
		}

		i.Add(&Import{
			Alias:         alias,
			PackagePath:   imp.PackagePath,
			ExplicitAlias: imp.ExplicitAlias || alias != imp.Alias,
		})
	}
}

type Import struct {
	Alias         string
	PackagePath   string
	ExplicitAlias bool
}

func parseImport(specs []*ast.ImportSpec) *Imports {
	imports := &Imports{
		Map:      container.NewOrderedMap[string, *Import](),
		AliasMap: map[string]string{},
	}

	for _, spec := range specs {
		i := Import{}

		i.PackagePath = mustStringLit(spec.Path)

		if spec.Name != nil {
			i.Alias = spec.Name.Name
			i.ExplicitAlias = true
		} else {
			i.Alias = i.PackagePath[strings.LastIndex(i.PackagePath, "/")+1:]
		}

		imports.Add(&i)
	}

	return imports
}
