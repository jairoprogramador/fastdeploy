package vos


// Version representa una versión semántica.
type Version struct {
	Major int
	// Minor es la versión menor (nuevas funcionalidades compatibles con versiones anteriores).
	Minor int
	// Patch es la versión de parche (correcciones de errores compatibles con versiones anteriores).
	Patch int
	// Raw es la representación en cadena de la versión completa (ej. "v1.2.3").
	Raw string
}

func (v *Version) String() string {
	return v.Raw
}