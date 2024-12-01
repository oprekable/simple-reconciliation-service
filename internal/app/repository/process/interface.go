package process

import (
	"context"
	"time"
)

type Repository interface {
	Pre(
		ctx context.Context,
		listBank []string,
		startDate time.Time,
		toDate time.Time,
	) (err error)

	ImportSystemTrx(ctx context.Context) (err error)
	Post(ctx context.Context) (err error)
	Close() (err error)
}
