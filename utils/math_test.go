package utils

import "testing"

func TestAdd(t *testing.T) {
	tests := []struct {
		name string
		a    int
		b    int
		want int
	}{
		{
			name: "positive numbers",
			a:    2,
			b:    3,
			want: 5,
		},
		{
			name: "negative numbers",
			a:    -2,
			b:    -3,
			want: -5,
		},
		{
			name: "positive and negative",
			a:    5,
			b:    -3,
			want: 2,
		},
		{
			name: "zero and positive",
			a:    0,
			b:    10,
			want: 10,
		},
		{
			name: "zero and negative",
			a:    0,
			b:    -7,
			want: -7,
		},
		{
			name: "both zeros",
			a:    0,
			b:    0,
			want: 0,
		},
		{
			name: "large numbers",
			a:    1000000,
			b:    2000000,
			want: 3000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Add(tt.a, tt.b); got != tt.want {
				t.Errorf("Add(%d, %d) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}