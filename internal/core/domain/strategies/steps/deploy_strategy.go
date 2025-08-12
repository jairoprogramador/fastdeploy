package steps

import "github.com/jairoprogramador/fastdeploy/internal/core/domain/context"

type DeployStrategy interface {
	ExecuteDeploy(context.Context) error
}
