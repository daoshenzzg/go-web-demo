package service

import (
	"context"
	"go-web-demo/app/demo-api/internal/conf"
	"go-web-demo/app/demo-api/internal/dao"
)

// Service struct
type Service struct {
	c   *conf.Config
	dao *dao.Dao
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
	}
	return
}

// Ping ping the resource.
func (s *Service) Ping(c context.Context) (err error) {
	// TODO
	return
}

// Close close the resource.
func (s *Service) Close() {
	s.dao.Close()
}
