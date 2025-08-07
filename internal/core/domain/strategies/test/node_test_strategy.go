package test

import "fmt"

type NodeTestStrategy struct{}

func NewNodeTestStrategy() TestStrategy {
	return &NodeTestStrategy{}
}

func (s *NodeTestStrategy) ExecuteTest() error {
	fmt.Println("  [Estrategia] Ejecutando pruebas para un proyecto Node.js (ej. npm test)")
	return nil
}