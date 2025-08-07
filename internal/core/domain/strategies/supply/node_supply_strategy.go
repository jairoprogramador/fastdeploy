package supply

import "fmt"

type NodeSupplyStrategy struct{}

func NewNodeSupplyStrategy() SupplyStrategy {
	return &NodeSupplyStrategy{}
}

func (s *NodeSupplyStrategy) ExecuteSupply() error {
	fmt.Println("  [Estrategia] Ejecutando supply para un proyecto Node.js (ej. infraestructura)")
	return nil
}
