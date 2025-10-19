package vos

import "errors"

type ProductID string

type Product struct {
	id           ProductID
	name         string
	description  string
	team         string
	organization string
}

func NewProduct(id ProductID, name, description, team, organization string) (Product, error) {
	if id == "" {
		return Product{}, errors.New("el ID del producto no puede estar vacío")
	}
	if name == "" {
		return Product{}, errors.New("el nombre del producto no puede estar vacío")
	}
	if organization == "" {
		return Product{}, errors.New("la organización del producto no puede estar vacía")
	}
	if team == "" {
		return Product{}, errors.New("el equipo del producto no puede estar vacío")
	}
	return Product{
		id:           id,
		name:         name,
		description:  description,
		team:         team,
		organization: organization,
	}, nil
}

func (p Product) ID() ProductID        { return p.id }
func (p Product) IdString() string     { return string(p.id) }
func (p Product) Name() string         { return p.name }
func (p Product) Description() string  { return p.description }
func (p Product) Team() string         { return p.team }
func (p Product) Organization() string { return p.organization }
