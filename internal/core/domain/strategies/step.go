package strategies

import "github.com/jairoprogramador/fastdeploy/internal/core/domain/context"

type Strategy interface {
	Execute(context.Context) error
}
