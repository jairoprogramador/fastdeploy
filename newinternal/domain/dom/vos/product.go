package vos

import "errors"

// ProductID es un hash que identifica unívocamente un producto.
type ProductID string

// Product agrupa los datos que definen un producto. Es un Objeto de Valor.
type Product struct {
	id           ProductID
	name         string
	description  string
	team         string
	organization string
}

func NewProduct(id ProductID, name, description, team, organization string) (*Product, error) {
	if id == "" {
		return nil, errors.New("el ID del producto no puede estar vacío")
	}
	if name == "" {
		return nil, errors.New("el nombre del producto no puede estar vacío")
	}
	if organization == "" {
		return nil, errors.New("la organización del producto no puede estar vacía")
	}
	if team == "" {
		return nil, errors.New("el equipo del producto no puede estar vacío")
	}
	return &Product{
		id:           id,
		name:         name,
		description:  description,
		team:         team,
		organization: organization,
	}, nil
}

// Getters para todos los campos...
func (p *Product) ID() ProductID        { return p.id }
func (p *Product) Name() string         { return p.name }
func (p *Product) Description() string  { return p.description }
func (p *Product) Team() string         { return p.team }
func (p *Product) Organization() string { return p.organization }
