package vos

import "time"

// Commit representa la información esencial de un commit de Git.
type Commit struct {
	// Hash es el SHA completo del commit.
	Hash string
	// Message es el mensaje de commit completo.
	Message string
	// Author es el autor del commit.
	Author string
	// Date es la fecha en que se realizó el commit.
	Date time.Time
}

// Version representa una versión semántica.
type Version struct {
	// Major es la versión mayor (cambios incompatibles).
	Major int
	// Minor es la versión menor (nuevas funcionalidades compatibles con versiones anteriores).
	Minor int
	// Patch es la versión de parche (correcciones de errores compatibles con versiones anteriores).
	Patch int
	// Raw es la representación en cadena de la versión completa (ej. "v1.2.3").
	Raw string
}
