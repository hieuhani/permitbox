package atomicity

import (
	"context"
)

type AtomicExecutor interface {
	Execute(parentCtx context.Context, executeFunc func(ctx context.Context) error) error
}
