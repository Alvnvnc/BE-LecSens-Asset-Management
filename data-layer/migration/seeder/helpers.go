package seeder

import "time"

// Helper functions used across seeder files

func floatPtr(f float64) *float64 {
	return &f
}

func intPtr(i int) *int {
	return &i
}

func int64Ptr(i int64) *int64 {
	return &i
}

func stringPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func safeFloat(f *float64) float64 {
	if f == nil {
		return 0.0
	}
	return *f
}

func safeString(s *string) string {
	if s == nil {
		return "unknown"
	}
	return *s
}
