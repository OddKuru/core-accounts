package log

import (
	"errors"
	"testing"
	"time"

	"github.com/OddKuru/accounts/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestField(t *testing.T) {
	allKeyName := "key"

	tc := []struct {
		name     string
		expected logger.Field
		actual   func() logger.Field
	}{
		{
			name: "Should correct string field",
			expected: logger.Field{
				Key:   allKeyName,
				Type:  logger.StringType,
				Value: "string value",
			},
			actual: func() logger.Field {
				return logger.String(allKeyName, "string value")
			},
		},
		{
			name: "Should correct int field",
			expected: logger.Field{
				Key:   allKeyName,
				Type:  logger.IntType,
				Value: 69,
			},
			actual: func() logger.Field {
				return logger.Int(allKeyName, 69)
			},
		},
		{
			name: "Should correct bool field",
			expected: logger.Field{
				Key:   allKeyName,
				Type:  logger.BoolType,
				Value: true,
			},
			actual: func() logger.Field {
				return logger.Bool(allKeyName, true)
			},
		},
		{
			name: "Should correct float64 field",
			expected: logger.Field{
				Key:   allKeyName,
				Type:  logger.Float64Type,
				Value: 69.69,
			},
			actual: func() logger.Field {
				return logger.Float64(allKeyName, 69.69)
			},
		},
		{
			name: "Should correct float32 field",
			expected: logger.Field{
				Key:   allKeyName,
				Type:  logger.Float32Type,
				Value: float32(69.69),
			},
			actual: func() logger.Field {
				return logger.Float32(allKeyName, 69.69)
			},
		},
		{
			name: "Should correct int8 field",
			expected: logger.Field{
				Key:   allKeyName,
				Type:  logger.Int8Type,
				Value: int8(69),
			},
			actual: func() logger.Field {
				return logger.Int8(allKeyName, 69)
			},
		},
		{
			name: "Should correct int16 field",
			expected: logger.Field{
				Key:   allKeyName,
				Type:  logger.Int16Type,
				Value: int16(69),
			},
			actual: func() logger.Field {
				return logger.Int16(allKeyName, 69)
			},
		},
		{
			name: "Should correct int32 field",
			expected: logger.Field{
				Key:   allKeyName,
				Type:  logger.Int32Type,
				Value: int32(69),
			},
			actual: func() logger.Field {
				return logger.Int32(allKeyName, 69)
			},
		},
		{
			name: "Should correct int64 field",
			expected: logger.Field{
				Key:   allKeyName,
				Type:  logger.Int64Type,
				Value: int64(69),
			},
			actual: func() logger.Field {
				return logger.Int64(allKeyName, 69)
			},
		},
		{
			name: "Should correct uint8 field",
			expected: logger.Field{
				Key:   allKeyName,
				Type:  logger.Uint8Type,
				Value: uint8(69),
			},
			actual: func() logger.Field {
				return logger.Uint8(allKeyName, 69)
			},
		},
		{
			name: "Should correct uint16 field",
			expected: logger.Field{
				Key:   allKeyName,
				Type:  logger.Uint16Type,
				Value: uint16(69),
			},
			actual: func() logger.Field {
				return logger.Uint16(allKeyName, 69)
			},
		},
		{
			name: "Should correct uint32 field",
			expected: logger.Field{
				Key:   allKeyName,
				Type:  logger.Uint32Type,
				Value: uint32(69),
			},
			actual: func() logger.Field {
				return logger.Uint32(allKeyName, 69)
			},
		},
		{
			name: "Should correct int64 field",
			expected: logger.Field{
				Key:   allKeyName,
				Type:  logger.Uint64Type,
				Value: uint64(69),
			},
			actual: func() logger.Field {
				return logger.Uint64(allKeyName, 69)
			},
		},
		{
			name: "Should correct time field",
			expected: logger.Field{
				Key:   allKeyName,
				Type:  logger.TimeType,
				Value: time.Unix(69, 0),
			},
			actual: func() logger.Field {
				return logger.Time(allKeyName, time.Unix(69, 0))
			},
		},
		{
			name: "Should correct Duration field",
			expected: logger.Field{
				Key:   allKeyName,
				Type:  logger.DurationType,
				Value: time.Duration(69) * time.Second,
			},
			actual: func() logger.Field {
				return logger.Duration(allKeyName, 69*time.Second)
			},
		},
		{
			name: "Should correct any field",
			expected: logger.Field{
				Key:   allKeyName,
				Type:  logger.AnyType,
				Value: []string{"kek", "lol"},
			},
			actual: func() logger.Field {
				return logger.Any(allKeyName, []string{"kek", "lol"})
			},
		},
		{
			name: "Should correct err field",
			expected: logger.Field{
				Key:   logger.ErrKey,
				Type:  logger.ErrorType,
				Value: errors.New("error 69"),
			},
			actual: func() logger.Field {
				return logger.Err(errors.New("error 69"))
			},
		},
		{
			name: "Should correct raw json field",
			expected: logger.Field{
				Key:   allKeyName,
				Type:  logger.RawJsonType,
				Value: []byte(`{"data": "message"}`),
			},
			actual: func() logger.Field {
				return logger.RawJson(allKeyName, []byte(`{"data": "message"}`))
			},
		},
		{
			name: "Should correct group field",
			expected: logger.Field{
				Key:  allKeyName,
				Type: logger.GroupType,
				Value: []logger.Field{
					logger.String("1", "mes"),
					logger.Int("2", 69),
				},
			},
			actual: func() logger.Field {
				return logger.Group(allKeyName, logger.String("1", "mes"), logger.Int("2", 69))
			},
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.actual())
		})
	}
}
