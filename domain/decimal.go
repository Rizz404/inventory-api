package domain

import (
	"encoding/json"
	"fmt"
)

// DecimalPrice is a custom type that always formats as decimal with 2 places
type DecimalPrice float64

// MarshalJSON formats the float64 as a decimal with 2 decimal places
func (d DecimalPrice) MarshalJSON() ([]byte, error) {
	// Format with 2 decimal places
	return json.Marshal(fmt.Sprintf("%.2f", float64(d)))
}

// UnmarshalJSON parses the decimal string back to float64
func (d *DecimalPrice) UnmarshalJSON(data []byte) error {
	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		// Try parsing as string if it fails as number
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		_, err := fmt.Sscanf(s, "%f", &f)
		if err != nil {
			return err
		}
	}
	*d = DecimalPrice(f)
	return nil
}

// Float64 returns the underlying float64 value
func (d DecimalPrice) Float64() float64 {
	return float64(d)
}

// DecimalValue is a custom type that always formats as decimal with 2 places (for nullable fields)
type DecimalValue struct {
	Value float64
	Valid bool
}

// MarshalJSON formats the float64 as a decimal with 2 decimal places or null
func (d DecimalValue) MarshalJSON() ([]byte, error) {
	if !d.Valid {
		return json.Marshal(nil)
	}
	// Format with 2 decimal places
	return json.Marshal(fmt.Sprintf("%.2f", d.Value))
}

// UnmarshalJSON parses the decimal string back to float64
func (d *DecimalValue) UnmarshalJSON(data []byte) error {
	// Check for null
	if string(data) == "null" {
		d.Valid = false
		return nil
	}

	var f float64
	if err := json.Unmarshal(data, &f); err != nil {
		// Try parsing as string if it fails as number
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		_, err := fmt.Sscanf(s, "%f", &f)
		if err != nil {
			return err
		}
	}
	d.Value = f
	d.Valid = true
	return nil
}
