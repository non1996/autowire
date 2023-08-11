package autowire

import (
	"github.com/non1996/go-jsonobj/function"
	"github.com/non1996/go-jsonobj/stream"
)

func GenerateAll(
	conf Config,
	root string,
) {
	dirs := traversalDir(root)

	dirStream := stream.Slice(dirs)
	packages := stream.
		MapS(dirStream, func(dir string) *Package { return parsePackage(conf, dir) }).
		Filter(function.NonNil[Package]).
		ToList()

	stream.Slice(packages).Foreach((*Package).generateFactories)
	stream.Slice(packages).Foreach((*Package).output)
}

func GenerateDir(
	conf Config,
	root string,
) {
	pkg := parsePackage(conf, root)
	if pkg == nil {
		return
	}
	pkg.generateFactories()
	pkg.output()
}
