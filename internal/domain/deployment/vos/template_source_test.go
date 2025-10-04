package vos

import "testing"

func TestNewTemplateSource(t *testing.T) {
	testCases := []struct {
		testName    string
		repoURL     string
		ref         string
		expectError bool
	}{
		{
			testName:    "Creacion valida con commit hash",
			repoURL:     "https://github.com/user/repo.git",
			ref:         "a1b2c3d4",
			expectError: false,
		},
		{
			testName:    "Creacion valida con tag",
			repoURL:     "https://github.com/user/repo.git",
			ref:         "v1.0.0",
			expectError: false,
		},
		{
			testName:    "Fallo por URL de repositorio vacia",
			repoURL:     "",
			ref:         "main",
			expectError: true,
		},
		{
			testName:    "Fallo por referencia vacia",
			repoURL:     "https://github.com/user/repo.git",
			ref:         "",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			source, err := NewTemplateSource(tc.repoURL, tc.ref)

			if tc.expectError {
				if err == nil {
					t.Errorf("Se esperaba un error, pero no se obtuvo ninguno")
				}
			} else {
				if err != nil {
					t.Errorf("No se esperaba un error, pero se obtuvo: %v", err)
				}
				if source.RepoURL() != tc.repoURL {
					t.Errorf("Se esperaba la URL '%s', pero se obtuvo '%s'", tc.repoURL, source.RepoURL())
				}
				if source.Ref() != tc.ref {
					t.Errorf("Se esperaba la referencia '%s', pero se obtuvo '%s'", tc.ref, source.Ref())
				}
			}
		})
	}
}
