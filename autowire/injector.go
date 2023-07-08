// generate by xxx

package autowire

type FieldInjector[T any] interface {
	Inject(*AppContext, *T)
}

type FieldInjectorImpl[C any, D any] struct {
	Qualifier string
	Require   bool
	InjectFn  func(*C, D)
}

func (f FieldInjectorImpl[C, D]) Inject(ctx *AppContext, comp *C) {
	var dep D

	if f.Qualifier != "" {
		dep = GetComponentByAlias[D](ctx, f.Qualifier, f.Require)
	} else {
		dep = GetComponent[D](ctx, f.Require)
	}

	f.InjectFn(comp, dep)
}
