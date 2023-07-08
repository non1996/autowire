package cli

import (
	"context"

	"github.com/non1996/go-jsonobj/stream"
	"golang.org/x/tools/go/packages"
)

func Load(
	ctx context.Context,
	dir string,
	env []string,
	tags string,
	patterns []string,
) ([]*packages.Package, []error) {
	cfg := &packages.Config{
		Context: ctx,
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedExportFile |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedTypesInfo |
			packages.NeedTypesSizes,
		Dir:        dir,
		Env:        env,
		BuildFlags: []string{},
		// TODO(light): Use ParseFile to skip function bodies and comments in indirect packages.
	}

	patterns = stream.Map(patterns, func(p string) string { return "pattern=" + p })

	pkgs, err := packages.Load(cfg, patterns...)
	if err != nil {
		return nil, []error{err}
	}

	var errs []error
	for _, p := range pkgs {
		for _, e := range p.Errors {
			errs = append(errs, e)
		}
	}

	if len(errs) > 0 {
		return nil, errs
	}

	return pkgs, nil
}
