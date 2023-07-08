package service

import (
	"github.com/non1996/autowire/example/client"
	"github.com/non1996/autowire/example/dal"
)

type TestService interface {
	Get() (int64, error)
	Set(int64) error
}

type TestServiceImpl struct {
	ADao dal.ADao
	BDao dal.BDao
	MQ   *client.MQ
	A    string

	m map[string]string
}

func (s *TestServiceImpl) Construct() {
	s.m = map[string]string{
		s.A: "123",
	}
}

func (s *TestServiceImpl) Get() (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (s *TestServiceImpl) Set(i int64) error {
	//TODO implement me
	panic("implement me")
}
