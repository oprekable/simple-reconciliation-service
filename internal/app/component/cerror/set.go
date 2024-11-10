package cerror

import (
	"simple-reconciliation-service/internal/app/err/core"

	"github.com/google/wire"
)

func ProvideErType(errType []core.ErrorType) ErType {
	return errType
}

var Set = wire.NewSet(
	ProvideErType,
	NewError,
)
