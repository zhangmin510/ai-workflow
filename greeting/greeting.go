package greeting

import "fmt"

// TimeOfDay represents the time of day for greeting purposes.
type TimeOfDay int

const (
	Morning TimeOfDay = iota
	Afternoon
	Evening
)

// greetingTemplate returns the greeting template based on formality and time.
// Uses functional style to map (isFormal, timeOfDay) → template string.
func greetingTemplate(isFormal bool, timeOfDay TimeOfDay) string {
	templates := map[TimeOfDay]map[bool]string{
		Morning: {
			true:  "Good morning, %s.",
			false: "Hey %s, good morning!",
		},
		Afternoon: {
			true:  "Good afternoon, %s.",
			false: "Hey %s, good afternoon!",
		},
		Evening: {
			true:  "Good evening, %s.",
			false: "Hey %s, good evening!",
		},
	}

	if timeTemplates, ok := templates[timeOfDay]; ok {
		if template, ok := timeTemplates[isFormal]; ok {
			return template
		}
	}

	// Default fallback
	if isFormal {
		return "Hello, %s."
	}
	return "Hey, %s!"
}

// GetGreeting generates a greeting message.
// name: the name to greet.
// isFormal: whether to use a formal tone.
// timeOfDay: the time of day (Morning, Afternoon, or Evening).
func GetGreeting(name string, isFormal bool, timeOfDay TimeOfDay) string {
	template := greetingTemplate(isFormal, timeOfDay)
	return fmt.Sprintf(template, name)
}
