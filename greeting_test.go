package greeting

import "testing"

func TestGetGreeting(t *testing.T) {
	tests := []struct {
		name      string
		isFormal  bool
		isMorning bool
		isEvening bool
		want      string
	}{
		{
			name:      "formal morning",
			isFormal:  true,
			isMorning: true,
			isEvening: false,
			want:      "Good morning, Alice.",
		},
		{
			name:      "formal evening",
			isFormal:  true,
			isMorning: false,
			isEvening: true,
			want:      "Good evening, Bob.",
		},
		{
			name:      "formal default",
			isFormal:  true,
			isMorning: false,
			isEvening: false,
			want:      "Hello, Charlie.",
		},
		{
			name:      "casual morning",
			isFormal:  false,
			isMorning: true,
			isEvening: false,
			want:      "Hey Alice, good morning!",
		},
		{
			name:      "casual evening",
			isFormal:  false,
			isMorning: false,
			isEvening: true,
			want:      "Hey Bob, good evening!",
		},
		{
			name:      "casual default",
			isFormal:  false,
			isMorning: false,
			isEvening: false,
			want:      "Hey, Charlie!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var name string
			if tt.isMorning {
				name = "Alice"
			} else if tt.isEvening {
				name = "Bob"
			} else {
				name = "Charlie"
			}

			got := GetGreeting(name, tt.isFormal, tt.isMorning, tt.isEvening)
			if got != tt.want {
				t.Errorf("GetGreeting() = %v, want %v", got, tt.want)
			}
		})
	}
}
