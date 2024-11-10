package service

import "github.com/google/wire"

func NewServices() *Services {
	return &Services{}
}

var Set = wire.NewSet(
	NewServices,
)
