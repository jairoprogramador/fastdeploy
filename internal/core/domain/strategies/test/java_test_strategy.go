package test

import "fmt"

type JavaTestStrategy struct{}

func NewJavaTestStrategy() TestStrategy {
	return &JavaTestStrategy{}
}

func (s *JavaTestStrategy) ExecuteTest() error {
	fmt.Println("  [Estrategia] Ejecutando pruebas para un proyecto Java (ej. mvn test)")
	return nil
}