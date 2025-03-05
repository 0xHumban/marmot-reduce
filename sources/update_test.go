package main

import (
	"testing"
)

func TestIsExit(t *testing.T) {
	tests := []struct {
		name     string
		message  Message
		expected bool
	}{
		{
			name: "Exit Message",
			message: Message{
				Type: String,
				Data: []byte("exit"),
			},
			expected: true,
		},
		{
			name: "Non-Exit Message",
			message: Message{
				Type: String,
				Data: []byte("hello"),
			},
			expected: false,
		},
		{
			name: "Non-String Type",
			message: Message{
				Type: 123, // Supposons que 123 n'est pas le type String
				Data: []byte("exit"),
			},
			expected: false,
		},
		{
			name: "Empty Data",
			message: Message{
				Type: String,
				Data: []byte(""),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.message.isExit()
			if result != tt.expected {
				t.Errorf("Expected %v, but got %v", tt.expected, result)
			}
		})
	}
}
