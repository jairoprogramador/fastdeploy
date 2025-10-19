package vos

import (
	deploymentvos "github.com/jairoprogramador/fastdeploy/internal/domain/deployment/vos"
	"testing"
)

func TestNewOutput(t *testing.T) {
	name := "project_name"
	value := "fastdeploy"

	outputDef, _ := deploymentvos.NewOutput(name, value)

	output := NewOutput(outputDef)

	if output.Name() != name {
		t.Errorf("Se esperaba el nombre '%s', pero se obtuvo '%s'", name, output.Name())
	}
	if output.Value() != value {
		t.Errorf("Se esperaba el valor '%s', pero se obtuvo '%s'", value, output.Value())
	}
}

func TestOutput_Equals(t *testing.T) {
	outputDef1, _ := deploymentvos.NewOutput("key1", "value1")
	v1 := NewOutput(outputDef1)
	outputDef2, _ := deploymentvos.NewOutput("key1", "value1")
	v2 := NewOutput(outputDef2)
	outputDef3, _ := deploymentvos.NewOutput("key2", "value1")
	v3 := NewOutput(outputDef3)
	outputDef4, _ := deploymentvos.NewOutput("key1", "value2")
	v4 := NewOutput(outputDef4)

	if !v1.Equals(v2) {
		t.Errorf("Se esperaba que v1 y v2 fueran iguales, pero no lo son")
	}

	if v1.Equals(v3) {
		t.Errorf("Se esperaba que v1 y v3 fueran diferentes, pero son iguales")
	}

	if v1.Equals(v4) {
		t.Errorf("Se esperaba que v1 y v4 fueran diferentes, pero son iguales")
	}
}
