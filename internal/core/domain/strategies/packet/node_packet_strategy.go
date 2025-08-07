package packet

import "fmt"

type NodePacketStrategy struct{}

func NewNodePacketStrategy() PacketStrategy {
	return &NodePacketStrategy{}
}

func (s *NodePacketStrategy) ExecutePacket() error {
	fmt.Println("  [Estrategia] Ejecutando package para un proyecto Node.js ")
	return nil
}
