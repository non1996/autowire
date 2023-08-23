package autowire

import (
	"go/ast"

	"github.com/non1996/go-jsonobj/container"

	"github.com/non1996/go-autowire/autowire-cli/annotation"
)

type GeneralFile struct {
	*Package
	AbsolutePath string
	Name         string
	Imports      Imports
	Structs      container.OrderedMap[string, *Struct]
	astFile      *ast.File
}

type AnnoFile struct {
	*Package
	Annotations []annotation.PrimaryAnnotation
}

type GenFile struct {
	*Package
	Imports      *Imports
	Components   []*Component
	Applications []*Application
	AbsolutePath string
	Output       []byte
}
