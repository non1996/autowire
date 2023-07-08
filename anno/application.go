package anno

type ApplicationAnnotation interface {
	isApplicationAnnotation()
}

type Application[_ any] []any

type ComponentScan struct {
	Value []string
}
