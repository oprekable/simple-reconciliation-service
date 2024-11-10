package err

import "simple-reconciliation-service/internal/app/err/core"

// RegisteredErrorType Register new errors here!
var RegisteredErrorType = []core.ErrorType{
	core.CErrInternal,
	core.CErrDBConn,
}
