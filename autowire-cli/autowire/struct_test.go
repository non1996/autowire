package autowire

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"testing"

	"github.com/non1996/go-jsonobj/container"

	"github.com/non1996/go-autowire/autowire-cli/internal/assert"
)

func TestPackageParseStruct(t *testing.T) {
	dir := "/Users/bytedance/goproj/autowire/example/dal"

	fSet = token.NewFileSet()
	packages, err := parser.ParseDir(fSet,
		dir,
		func(info fs.FileInfo) bool {
			return info.Name() != "autowire.go" &&
				info.Name() != "autowire_gen.go"
		},
		parser.DeclarationErrors,
	)
	if err != nil {
		panic(err)
	}
	if len(packages) == 0 {
		return
	}
	assert.Assert(len(packages) == 1)

	astPkg := packages[container.MapKeys(packages)[0]]

	p := NewPackage(Config{}, dir, astPkg)

	p.parseStructs()

	fmt.Println(p.Structs)
}
