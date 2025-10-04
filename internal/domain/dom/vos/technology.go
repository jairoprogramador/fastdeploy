package vos

import "errors"

// Technology agrupa los datos que definen una tecnología. Es un Objeto de Valor.
type Technology struct {
	type_technology string
	solution        string
	stack           string
	infrastructure  string
}

func NewTechnology(type_technology, solution, stack, infrastructure string) (*Technology, error) {
	if type_technology == "" {
		return nil, errors.New("el tipo de tecnología no puede estar vacío")
	}
	if solution == "" {
		return nil, errors.New("la solución no puede estar vacía")
	}
	if stack == "" {
		return nil, errors.New("la stack no puede estar vacía")
	}
	if infrastructure == "" {
		return nil, errors.New("la infraestructura no puede estar vacía")
	}
	return &Technology{
		type_technology: type_technology,
		solution:        solution,
		stack:           stack,
		infrastructure:  infrastructure,
	}, nil
}

// Getters para todos los campos...
func (t *Technology) TypeTechnology() string { return t.type_technology }
func (t *Technology) Solution() string       { return t.solution }
func (t *Technology) Stack() string          { return t.stack }
func (t *Technology) Infrastructure() string { return t.infrastructure }