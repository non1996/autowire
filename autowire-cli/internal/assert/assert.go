package assert

import (
	"fmt"

	"github.com/non1996/go-jsonobj/container"
)

func Assert(cond bool, msg ...string) {
	if !cond {
		panic(fmt.Errorf("assert failed, %s", container.SliceGetFirst(msg)))
	}
}

func AssertFunc(cond func() bool, msg ...any) {
	if !cond() {
		panic(fmt.Errorf("assert failed, %s", container.SliceGetFirst(msg)))
	}
}
