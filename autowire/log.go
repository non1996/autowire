package autowire

type LogHandler interface {
	Debug(format string, args ...any)
	Info(format string, args ...any)
	Error(format string, args ...any)
}

type defaultLogHandler struct {
}

func (d defaultLogHandler) Debug(format string, args ...any) {
}

func (d defaultLogHandler) Info(format string, args ...any) {
}

func (d defaultLogHandler) Error(format string, args ...any) {
}
