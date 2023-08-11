package autowire

import (
	"go/ast"
	"os"
	"path"
	"strings"

	"github.com/modern-go/reflect2"
)

func traversalDir(root string) (res []string) {
	entries, err := os.ReadDir(root)
	if err != nil {
		panic(err)
	}

	res = append(res, root)
	for _, entry := range entries {
		if entry.IsDir() {
			res = append(res, traversalDir(path.Join(root, entry.Name()))...)
		}
	}

	return res
}

func mustStringLit(expr ast.Expr, defaultValue ...string) string {
	if reflect2.IsNil(expr) {
		if len(defaultValue) == 0 {
			panic("")
		}
		return defaultValue[0]
	}

	return strings.Trim(expr.(*ast.BasicLit).Value, "\\\"|`")
}

func mustIdent(expr ast.Expr) *ast.Ident {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		panic("")
	}

	return ident
}
