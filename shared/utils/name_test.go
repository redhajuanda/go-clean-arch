package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeName(t *testing.T) {
	testTable := []struct {
		input string
		want  []string
	}{
		{
			input: "Redha Arifan Juanda  ",
			want:  []string{"Redha", "Arifan Juanda"},
		},
		{
			input: "Barrack Obama",
			want:  []string{"Barrack", "Obama"},
		},
		{
			input: "Husein",
			want:  []string{"Husein", ""},
		},
		{
			input: "",
			want:  []string{"", ""},
		},
		{
			input: "Cinta Rangga AADC Suara",
			want:  []string{"Cinta", "Rangga AADC Suara"},
		},
		{
			input: "Kembang  Desa",
			want:  []string{"Kembang", "Desa"},
		},
		{
			input: "  Bunga   Mawar",
			want:  []string{"Bunga", "Mawar"},
		},
	}

	for _, test := range testTable {
		gotFirstName, gotLastName := NormalizeName(test.input)
		assert.Equal(t, test.want[0], gotFirstName)
		assert.Equal(t, test.want[1], gotLastName)
	}
}
