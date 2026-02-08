package logger

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestField(t *testing.T) {
	allKeyName := "key"

	tc := []struct {
		name     string
		expected Field
		actual   func() Field
	}{
		{
			name: "Should correct string field",
			expected: Field{
				Key:   allKeyName,
				Type:  StringType,
				Value: "string value",
			},
			actual: func() Field {
				return String(allKeyName, "string value")
			},
		},
		{
			name: "Should correct int field",
			expected: Field{
				Key:   allKeyName,
				Type:  IntType,
				Value: 69,
			},
			actual: func() Field {
				return Int(allKeyName, 69)
			},
		},
		{
			name: "Should correct bool field",
			expected: Field{
				Key:   allKeyName,
				Type:  BoolType,
				Value: true,
			},
			actual: func() Field {
				return Bool(allKeyName, true)
			},
		},
		{
			name: "Should correct float64 field",
			expected: Field{
				Key:   allKeyName,
				Type:  Float64Type,
				Value: 69.69,
			},
			actual: func() Field {
				return Float64(allKeyName, 69.69)
			},
		},
		{
			name: "Should correct float32 field",
			expected: Field{
				Key:   allKeyName,
				Type:  Float32Type,
				Value: float32(69.69),
			},
			actual: func() Field {
				return Float32(allKeyName, 69.69)
			},
		},
		{
			name: "Should correct int8 field",
			expected: Field{
				Key:   allKeyName,
				Type:  Int8Type,
				Value: int8(69),
			},
			actual: func() Field {
				return Int8(allKeyName, 69)
			},
		},
		{
			name: "Should correct int16 field",
			expected: Field{
				Key:   allKeyName,
				Type:  Int16Type,
				Value: int16(69),
			},
			actual: func() Field {
				return Int16(allKeyName, 69)
			},
		},
		{
			name: "Should correct int32 field",
			expected: Field{
				Key:   allKeyName,
				Type:  Int32Type,
				Value: int32(69),
			},
			actual: func() Field {
				return Int32(allKeyName, 69)
			},
		},
		{
			name: "Should correct int64 field",
			expected: Field{
				Key:   allKeyName,
				Type:  Int64Type,
				Value: int64(69),
			},
			actual: func() Field {
				return Int64(allKeyName, 69)
			},
		},
		{
			name: "Should correct uint8 field",
			expected: Field{
				Key:   allKeyName,
				Type:  Uint8Type,
				Value: uint8(69),
			},
			actual: func() Field {
				return Uint8(allKeyName, 69)
			},
		},
		{
			name: "Should correct uint16 field",
			expected: Field{
				Key:   allKeyName,
				Type:  Uint16Type,
				Value: uint16(69),
			},
			actual: func() Field {
				return Uint16(allKeyName, 69)
			},
		},
		{
			name: "Should correct uint32 field",
			expected: Field{
				Key:   allKeyName,
				Type:  Uint32Type,
				Value: uint32(69),
			},
			actual: func() Field {
				return Uint32(allKeyName, 69)
			},
		},
		{
			name: "Should correct int64 field",
			expected: Field{
				Key:   allKeyName,
				Type:  Uint64Type,
				Value: uint64(69),
			},
			actual: func() Field {
				return Uint64(allKeyName, 69)
			},
		},
		{
			name: "Should correct time field",
			expected: Field{
				Key:   allKeyName,
				Type:  TimeType,
				Value: time.Unix(69, 0),
			},
			actual: func() Field {
				return Time(allKeyName, time.Unix(69, 0))
			},
		},
		{
			name: "Should correct Duration field",
			expected: Field{
				Key:   allKeyName,
				Type:  DurationType,
				Value: time.Duration(69) * time.Second,
			},
			actual: func() Field {
				return Duration(allKeyName, 69*time.Second)
			},
		},
		{
			name: "Should correct any field",
			expected: Field{
				Key:   allKeyName,
				Type:  AnyType,
				Value: []string{"kek", "lol"},
			},
			actual: func() Field {
				return Any(allKeyName, []string{"kek", "lol"})
			},
		},
		{
			name: "Should correct err field",
			expected: Field{
				Key:   ErrKey,
				Type:  ErrorType,
				Value: errors.New("error 69"),
			},
			actual: func() Field {
				return Err(errors.New("error 69"))
			},
		},
		{
			name: "Should correct raw json field",
			expected: Field{
				Key:   allKeyName,
				Type:  RawJsonType,
				Value: []byte(`{"data": "message"}`),
			},
			actual: func() Field {
				return RawJson(allKeyName, []byte(`{"data": "message"}`))
			},
		},
		{
			name: "Should correct group field",
			expected: Field{
				Key:  allKeyName,
				Type: GroupType,
				Value: []Field{
					String("1", "mes"),
					Int("2", 69),
				},
			},
			actual: func() Field {
				return Group(allKeyName, String("1", "mes"), Int("2", 69))
			},
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.actual())
		})
	}
}
