package autowire

import (
	"github.com/non1996/go-jsonobj/function"
	"github.com/non1996/go-jsonobj/stream"
)

func Generate(
	conf Config,
) {
	dirs := traversalDir(conf.Root)

	dirStream := stream.Slice(dirs)

	genContext.packages = stream.
		MapS(dirStream, func(dir string) *Package { return parsePackage(conf, dir) }).
		Filter(function.NonNil[Package]).
		ToList()

	stream.Slice(genContext.packages).Foreach((*Package).evaluate)
	stream.Slice(genContext.packages).Foreach((*Package).format)
	stream.Slice(genContext.packages).Foreach((*Package).output)
}
