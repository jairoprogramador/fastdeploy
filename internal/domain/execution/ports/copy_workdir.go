package ports

import (
	"context"
)

type CopyWorkdir interface {
	Copy(ctx context.Context, source, destination string) error
}