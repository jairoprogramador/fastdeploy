package vos

import "testing"

func TestNewFingerprint(t *testing.T) {
	testCases := []struct {
		name    string
		value   string
		want    Fingerprint
		wantErr bool
	}{
		{
			name:    "debería crear un fingerprint válido con un valor no vacío",
			value:   "some-valid-fingerprint",
			want:    Fingerprint{value: "some-valid-fingerprint"},
			wantErr: false,
		},
		{
			name:    "debería devolver un error si el valor está vacío",
			value:   "",
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := NewFingerprint(tc.value)

			if (err != nil) != tc.wantErr {
				t.Errorf("NewFingerprint() error = %v, wantErr %v", err, tc.wantErr)
				return
			}

			if !tc.wantErr && got != tc.want {
				t.Errorf("NewFingerprint() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestFingerprint_String(t *testing.T) {
	valueFingerprint := "my-fingerprint"

	fp, err := NewFingerprint(valueFingerprint)
	if err != nil {
		t.Fatalf("fallo en la configuración de la prueba: no se pudo crear el fingerprint: %v", err)
	}

	testCases := []struct {
		name string
		fingerprint Fingerprint
		want string
	}{
		{
			name: "debería devolver el string para un fingerprint válido",
			fingerprint:    fp,
			want: valueFingerprint,
		},
		{
			name: "debería devolver un string vacío para un fingerprint de valor cero",
			fingerprint:    Fingerprint{},
			want: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.fingerprint.String(); got != tc.want {
				t.Errorf("String() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestFingerprint_Equals(t *testing.T) {
	fp1, _ := NewFingerprint("same-fingerprint")
	fp1Clone, _ := NewFingerprint("same-fingerprint")
	fp2, _ := NewFingerprint("different-fingerprint")

	testCases := []struct {
		name string
		fingerprint1   Fingerprint
		fingerprint2   Fingerprint
		want bool
	}{
		{
			name: "debería devolver verdadero para dos fingerprints con el mismo valor",
			fingerprint1:   fp1,
			fingerprint2:   fp1Clone,
			want: true,
		},
		{
			name: "debería devolver falso para dos fingerprints con valores diferentes",
			fingerprint1:   fp1,
			fingerprint2:   fp2,
			want: false,
		},
		{
			name: "debería devolver falso al comparar un fingerprint válido con uno de valor cero",
			fingerprint1:   fp1,
			fingerprint2:   Fingerprint{},
			want: false,
		},
		{
			name: "debería devolver verdadero al comparar dos fingerprints de valor cero",
			fingerprint1:   Fingerprint{},
			fingerprint2:   Fingerprint{},
			want: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.fingerprint1.Equals(tc.fingerprint2); got != tc.want {
				t.Errorf("Equals() = %v, want %v", got, tc.want)
			}
		})
	}
}
