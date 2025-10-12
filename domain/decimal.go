package domain

import (
	"encoding/json"
	"fmt"
)

// Decimal2 is a custom float64 type that always marshals to JSON with exactly 2 decimal places as a number
// Example: 3211.0 -> 3211.00, 38.095238095238095 -> 38.10
type Decimal2 float64

// MarshalJSON formats the float64 as a number with exactly 2 decimal places
func (d Decimal2) MarshalJSON() ([]byte, error) {
	// Format to 2 decimal places using fmt.Sprintf to ensure .00 is always present
	formatted := fmt.Sprintf("%.2f", float64(d))
	return []byte(formatted), nil
}

// UnmarshalJSON parses a JSON number to Decimal2
func (d *Decimal2) UnmarshalJSON(data []byte) error {
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	*d = Decimal2(f)
	return nil
}

// Float64 returns the underlying float64 value
func (d Decimal2) Float64() float64 {
	return float64(d)
}

// NewDecimal2 creates a new Decimal2 from a float64
func NewDecimal2(value float64) Decimal2 {
	return Decimal2(value)
}

// NullableDecimal2 is a nullable version of Decimal2
type NullableDecimal2 struct {
	Value Decimal2
	Valid bool
}

// MarshalJSON formats the float64 as a number with 2 decimal places or null
// Uses pointer receiver to ensure it works with both pointer and value types
func (d *NullableDecimal2) MarshalJSON() ([]byte, error) {
	if d == nil || !d.Valid {
		return []byte("null"), nil
	}
	// Use fmt.Sprintf to ensure .00 is always present
	formatted := fmt.Sprintf("%.2f", float64(d.Value))
	return []byte(formatted), nil
}

// UnmarshalJSON parses a JSON number or null to NullableDecimal2
func (d *NullableDecimal2) UnmarshalJSON(data []byte) error {
	// Check for null
	if string(data) == "null" {
		d.Valid = false
		return nil
	}

	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		return err
	}
	d.Value = Decimal2(f)
	d.Valid = true
	return nil
}

// Float64 returns the underlying float64 value if valid
func (d NullableDecimal2) Float64() (float64, bool) {
	if !d.Valid {
		return 0, false
	}
	return d.Value.Float64(), true
}

// NewNullableDecimal2 creates a new NullableDecimal2 from a *float64
func NewNullableDecimal2(value *float64) *NullableDecimal2 {
	if value == nil {
		return &NullableDecimal2{Valid: false}
	}
	return &NullableDecimal2{Value: Decimal2(*value), Valid: true}
}
