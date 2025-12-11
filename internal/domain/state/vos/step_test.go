package vos

import "testing"

func TestNewStep(t *testing.T) {
	testCases := []struct {
		name    string
		value   string
		want    Step
		wantErr bool
	}{
		{
			name:    "debería crear un step válido para test",
			value:   StepTest,
			want:    Step{value: StepTest},
			wantErr: false,
		},
		{
			name:    "debería crear un step válido para supply",
			value:   StepSupply,
			want:    Step{value: StepSupply},
			wantErr: false,
		},
		{
			name:    "debería crear un step válido para package",
			value:   StepPackage,
			want:    Step{value: StepPackage},
			wantErr: false,
		},
		{
			name:    "debería crear un step válido para deploy",
			value:   StepDeploy,
			want:    Step{value: StepDeploy},
			wantErr: false,
		},
		{
			name:    "debería devolver un error para un valor vacío",
			value:   "",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewStep(tc.value)

			if (err != nil) != tc.wantErr {
				t.Errorf("NewStep() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr && got != tc.want {
				t.Errorf("NewStep() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestStep_String(t *testing.T) {
	testCases := []struct {
		step Step
		want string
	}{
		{Step{value: StepTest}, "test"},
		{Step{value: StepSupply}, "supply"},
		{Step{value: StepPackage}, "package"},
		{Step{value: StepDeploy}, "deploy"},
		{Step{value: ""}, ""}, // Probando el valor cero
	}

	for _, tc := range testCases {
		if got := tc.step.String(); got != tc.want {
			t.Errorf("String() para %q = %q, want %q", tc.step, got, tc.want)
		}
	}
}
