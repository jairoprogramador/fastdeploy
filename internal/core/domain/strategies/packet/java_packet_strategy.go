package packet

import "fmt"

type JavaPacketStrategy struct{}

func NewJavaPacketStrategy() PacketStrategy {
	return &JavaPacketStrategy{}
}

func (s *JavaPacketStrategy) ExecutePacket() error {
	fmt.Println("  [Estrategia] Ejecutando package para un proyecto Java")
	return nil
}
