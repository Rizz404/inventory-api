package utils

// Ptr returns a pointer to the given value (generic version)
func Ptr[T any](v T) *T {
	return &v
}

// StringPtr returns a pointer to string
func StringPtr(s string) *string {
	return &s
}

// BoolPtr returns a pointer to bool
func BoolPtr(b bool) *bool {
	return &b
}

// Float64Ptr returns a pointer to float64
func Float64Ptr(f float64) *float64 {
	return &f
}

// IntPtr returns a pointer to int
func IntPtr(i int) *int {
	return &i
}

// Int64Ptr returns a pointer to int64
func Int64Ptr(i int64) *int64 {
	return &i
}
