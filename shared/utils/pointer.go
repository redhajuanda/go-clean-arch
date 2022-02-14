package utils

import "time"

// StringPtr returns a pointer to the passed string.
func StringPtr(s string) *string {
	return &s
}

// IntPtr returns a pointer to the passed int.
func IntPtr(i int) *int {
	return &i
}

// MapStringInterfacePtr returns a pointer to the passed map string interface.
func MapStringInterfacePtr(m map[string]interface{}) *map[string]interface{} {
	return &m
}

// TimePtr returns a pointer to the passed time.
func TimePtr(t time.Time) *time.Time {
	return &t
}

func PtrStr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func PtrInt(s *int) int {
	if s == nil {
		return 0
	}
	return *s
}

func PtrFloat64(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}
