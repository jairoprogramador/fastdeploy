package steps

import "github.com/jairoprogramador/fastdeploy/internal/core/domain/context"

type TestStrategy interface {
	ExecuteTest(context.Context) error
}
