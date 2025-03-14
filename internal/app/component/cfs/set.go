package cfs

import (
	"github.com/google/wire"
	"github.com/spf13/afero"
)

type FSType struct {
	LocalStorageFs afero.Fs
}

func ProviderCFs(fsType FSType) *Fs {
	return NewFs(
		fsType.LocalStorageFs,
	)
}

var Set = wire.NewSet(
	ProviderCFs,
)
