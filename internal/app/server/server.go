package server

import (
	"reflect"
	"simple-reconciliation-service/internal/app/server/cli"
	"simple-reconciliation-service/internal/pkg/utils/atexit"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	Cli *cli.Cli
}

func NewServer(
	cli *cli.Cli,
) *Server {
	return &Server{
		Cli: cli,
	}
}

func (s *Server) Run(eg *errgroup.Group) {
	v := reflect.ValueOf(*s)
	for i := 0; i < v.NumField(); i++ {
		if s, ok := v.Field(i).Interface().(IServer); ok {
			atexit.Add(s.Shutdown)
			s.Start(eg)
		}
	}
}
