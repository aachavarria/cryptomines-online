package utils

import "testing"

func TestContainsProfanity(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected bool
	}{
		// Test cases with profanity
		{
			name:     "contains fuck",
			message:  "what the fuck is this",
			expected: true,
		},
		{
			name:     "contains shit",
			message:  "this is shit",
			expected: true,
		},
		{
			name:     "contains bitch",
			message:  "stop being a bitch",
			expected: true,
		},
		{
			name:     "contains asshole",
			message:  "you're such an asshole",
			expected: true,
		},
		{
			name:     "contains damn",
			message:  "damn it all",
			expected: true,
		},
		{
			name:     "contains crap",
			message:  "this is crap",
			expected: true,
		},
		{
			name:     "contains profanity with caps",
			message:  "FUCK THIS SHIT",
			expected: true,
		},
		{
			name:     "contains profanity mixed case",
			message:  "FuCk ThIs",
			expected: true,
		},
		{
			name:     "contains profanity in sentence",
			message:  "I can't believe this fucking happened",
			expected: true,
		},
		{
			name:     "contains multiple profanity words",
			message:  "this shit is fucking terrible",
			expected: true,
		},
		// Test cases without profanity
		{
			name:     "clean message",
			message:  "Hello world, how are you?",
			expected: false,
		},
		{
			name:     "game related message",
			message:  "Let's attack that planet!",
			expected: false,
		},
		{
			name:     "tactical message",
			message:  "Need backup at sector 7",
			expected: false,
		},
		{
			name:     "friendly greeting",
			message:  "GG everyone, great game!",
			expected: false,
		},
		{
			name:     "empty string",
			message:  "",
			expected: false,
		},
		{
			name:     "numbers only",
			message:  "12345",
			expected: false,
		},
		{
			name:     "special characters",
			message:  "!@#$%^&*()",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ContainsProfanity(tt.message)
			if result != tt.expected {
				t.Errorf("ContainsProfanity(%q) = %v, expected %v", tt.message, result, tt.expected)
			}
		})
	}
}

func TestProfanityFilterCaseInsensitive(t *testing.T) {
	testMessages := []string{
		"fuck",
		"FUCK",
		"FuCk",
		"fUcK",
		"FuCK",
	}

	for _, msg := range testMessages {
		if !ContainsProfanity(msg) {
			t.Errorf("ContainsProfanity(%q) should return true for case variations", msg)
		}
	}
}

func TestProfanityFilterComprehensive(t *testing.T) {
	// Test that we have a comprehensive list (at least 50 words)
	if len(profanityWords) < 50 {
		t.Errorf("Expected at least 50 profanity words in the list, got %d", len(profanityWords))
	}

	// Test a sample of common profanity words
	commonProfanity := []string{
		"fuck", "shit", "bitch", "asshole", "bastard",
		"damn", "crap", "ass", "dick", "pussy",
	}

	for _, word := range commonProfanity {
		if !ContainsProfanity(word) {
			t.Errorf("Common profanity word %q should be blocked", word)
		}
	}
}
