package main

import (
	"slices"
	"testing"
)

func TestStatements(t *testing.T) {
	tests := []struct {
		name string
		sql  string
		want []string
	}{
		{"um por linha", "SELECT * FROM a;\nSELECT * FROM b;", []string{"SELECT * FROM a;", "SELECT * FROM b;"}},
		{"várias linhas", "CREATE TABLE p (\n  id INT PRIMARY KEY,\n  nome TEXT\n);", []string{"CREATE TABLE p ( id INT PRIMARY KEY, nome TEXT );"}},
		{"comentários e linhas vazias", "-- oi\n\nSELECT 1 FROM a;\n  -- tchau\n", []string{"SELECT 1 FROM a;"}},
		{"último sem ponto e vírgula", "SELECT * FROM a;\nSELECT *\nFROM b", []string{"SELECT * FROM a;", "SELECT * FROM b"}},
		{"vazio", "\n  \n", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := statements(tt.sql); !slices.Equal(got, tt.want) {
				t.Errorf("statements(%q) = %q, queria %q", tt.sql, got, tt.want)
			}
		})
	}
}
