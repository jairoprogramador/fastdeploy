package steps

import "github.com/jairoprogramador/fastdeploy/internal/core/domain/context"

type PacketStrategy interface {
	ExecutePacket(context.Context) error
}
