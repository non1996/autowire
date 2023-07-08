package handler

import (
	"context"

	"github.com/non1996/autowire/example/service"
)

type TestHandler struct {
	TestService service.TestService
}

func (h *TestHandler) Register() {
	register(h.Handle)
}

func (h *TestHandler) Handle(ctx context.Context) error {
	a, err := h.TestService.Get()
	if err != nil {
		return err
	}
	err = h.TestService.Set(a + 1)
	if err != nil {
		return err
	}

	return nil
}

func register(_ any) {

}
