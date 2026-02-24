package utils

import (
	"testing"
	"time"
)

func TestTimeAgo(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		time     time.Time
		expected string
	}{
		{
			name:     "just now - 30 seconds ago",
			time:     now.Add(-30 * time.Second),
			expected: "just now",
		},
		{
			name:     "1 minute ago",
			time:     now.Add(-1 * time.Minute),
			expected: "1 minute ago",
		},
		{
			name:     "5 minutes ago",
			time:     now.Add(-5 * time.Minute),
			expected: "5 minutes ago",
		},
		{
			name:     "59 minutes ago",
			time:     now.Add(-59 * time.Minute),
			expected: "59 minutes ago",
		},
		{
			name:     "1 hour ago",
			time:     now.Add(-1 * time.Hour),
			expected: "1 hour ago",
		},
		{
			name:     "3 hours ago",
			time:     now.Add(-3 * time.Hour),
			expected: "3 hours ago",
		},
		{
			name:     "23 hours ago",
			time:     now.Add(-23 * time.Hour),
			expected: "23 hours ago",
		},
		{
			name:     "1 day ago",
			time:     now.Add(-24 * time.Hour),
			expected: "1 day ago",
		},
		{
			name:     "7 days ago",
			time:     now.Add(-7 * 24 * time.Hour),
			expected: "7 days ago",
		},
		{
			name:     "29 days ago",
			time:     now.Add(-29 * 24 * time.Hour),
			expected: "29 days ago",
		},
		{
			name:     "1 month ago",
			time:     now.Add(-30 * 24 * time.Hour),
			expected: "1 month ago",
		},
		{
			name:     "3 months ago",
			time:     now.Add(-90 * 24 * time.Hour),
			expected: "3 months ago",
		},
		{
			name:     "11 months ago",
			time:     now.Add(-330 * 24 * time.Hour),
			expected: "11 months ago",
		},
		{
			name:     "1 year ago",
			time:     now.Add(-365 * 24 * time.Hour),
			expected: "1 year ago",
		},
		{
			name:     "2 years ago",
			time:     now.Add(-730 * 24 * time.Hour),
			expected: "2 years ago",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TimeAgo(tt.time)
			if got != tt.expected {
				t.Errorf("TimeAgo(%v) = %q, want %q", tt.time, got, tt.expected)
			}
		})
	}
}

func TestTimeAgo_EdgeCases(t *testing.T) {
	t.Run("future time returns just now", func(t *testing.T) {
		future := time.Now().Add(5 * time.Minute)
		got := TimeAgo(future)
		if got != "just now" {
			t.Errorf("TimeAgo(future) = %q, want %q", got, "just now")
		}
	})

	t.Run("exactly 1 minute", func(t *testing.T) {
		oneMinuteAgo := time.Now().Add(-1 * time.Minute)
		got := TimeAgo(oneMinuteAgo)
		if got != "1 minute ago" {
			t.Errorf("TimeAgo(1 minute ago) = %q, want %q", got, "1 minute ago")
		}
	})

	t.Run("exactly 1 hour", func(t *testing.T) {
		oneHourAgo := time.Now().Add(-1 * time.Hour)
		got := TimeAgo(oneHourAgo)
		if got != "1 hour ago" {
			t.Errorf("TimeAgo(1 hour ago) = %q, want %q", got, "1 hour ago")
		}
	})
}
