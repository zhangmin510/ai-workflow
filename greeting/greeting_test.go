package greeting

import "testing"

func TestGetGreeting(t *testing.T) {
	tests := []struct {
		name      string
		isFormal  bool
		timeOfDay TimeOfDay
		personName string
		want      string
	}{
		{
			name:      "formal morning",
			isFormal:  true,
			timeOfDay: Morning,
			personName: "Alice",
			want:      "Good morning, Alice.",
		},
		{
			name:      "formal afternoon",
			isFormal:  true,
			timeOfDay: Afternoon,
			personName: "Bob",
			want:      "Hello, Bob.",
		},
		{
			name:      "formal evening",
			isFormal:  true,
			timeOfDay: Evening,
			personName: "Charlie",
			want:      "Good evening, Charlie.",
		},
		{
			name:      "casual morning",
			isFormal:  false,
			timeOfDay: Morning,
			personName: "Alice",
			want:      "Hey Alice, good morning!",
		},
		{
			name:      "casual afternoon",
			isFormal:  false,
			timeOfDay: Afternoon,
			personName: "Bob",
			want:      "Hey, Bob!",
		},
		{
			name:      "casual evening",
			isFormal:  false,
			timeOfDay: Evening,
			personName: "Charlie",
			want:      "Hey Charlie, good evening!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetGreeting(tt.personName, tt.isFormal, tt.timeOfDay)
			if got != tt.want {
				t.Errorf("GetGreeting() = %v, want %v", got, tt.want)
			}
		})
	}
}
