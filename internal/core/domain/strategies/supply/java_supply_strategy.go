package supply

import "fmt"

type JavaSupplyStrategy struct{}

func NewJavaSupplyStrategy() SupplyStrategy {
	return &JavaSupplyStrategy{}
}

func (s *JavaSupplyStrategy) ExecuteSupply() error {
	fmt.Println("  [Estrategia] Ejecutando supply para un proyecto Java (ej. infraestructura)")
	return nil
}
