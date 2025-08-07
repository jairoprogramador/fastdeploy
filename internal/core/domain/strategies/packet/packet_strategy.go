package packet

type PacketStrategy interface {
	ExecutePacket() error
}
