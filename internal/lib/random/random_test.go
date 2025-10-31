package random

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRandomString(t *testing.T) {
	tests := []struct {
		name string
		size int
	}{
		{
			name: "size = 1",
			size: 1,
		},
		{
			name: "size = 5",
			size: 5,
		},
		{
			name: "size = 10",
			size: 10,
		},
		{
			name: "size = 20",
			size: 20,
		},
		{
			name: "size = 30",
			size: 30,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str1 := NewRandomString(tt.size)
			str2 := NewRandomString(tt.size)

			fmt.Println(str1)
			fmt.Println(str2)

			assert.Len(t, str1, tt.size)
			assert.Len(t, str2, tt.size)

			// проверяем что две сгенерированные строки отличаются

			assert.NotEqual(t, str1, str2)
		})
	}
}
