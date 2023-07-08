package anno

type Configuration[_ any] []any

type Bean struct {
	Name       string
	Method     any
	InitMethod any
}
