package logger

import "time"

func String(key string, value string) Field {
	return Field{Key: key, Type: StringType, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Type: IntType, Value: value}
}

func Any(key string, value interface{}) Field {
	return Field{Key: key, Type: AnyType, Value: value}
}

func Err(err error) Field {
	return Field{Key: ErrKey, Type: ErrorType, Value: err}
}

func Duration(key string, value time.Duration) Field {
	return Field{Key: key, Type: DurationType, Value: value}
}

func Bool(key string, value bool) Field {
	return Field{Key: key, Type: BoolType, Value: value}
}

func Time(key string, value time.Time) Field {
	return Field{Key: key, Type: TimeType, Value: value}
}

func Float32(key string, value float32) Field {
	return Field{Key: key, Type: Float32Type, Value: value}
}

func Float64(key string, value float64) Field {
	return Field{Key: key, Type: Float64Type, Value: value}
}

func Int8(key string, value int8) Field {
	return Field{Key: key, Type: Int8Type, Value: value}
}

func Int16(key string, value int16) Field {
	return Field{Key: key, Type: Int16Type, Value: value}
}

func Int32(key string, value int32) Field {
	return Field{Key: key, Type: Int32Type, Value: value}
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Type: Int64Type, Value: value}
}

func Uint8(key string, value uint8) Field {
	return Field{Key: key, Type: Uint8Type, Value: value}
}

func Uint16(key string, value uint16) Field {
	return Field{Key: key, Type: Uint16Type, Value: value}
}

func Uint32(key string, value uint32) Field {
	return Field{Key: key, Type: Uint32Type, Value: value}
}

func Uint64(key string, value uint64) Field {
	return Field{Key: key, Type: Uint64Type, Value: value}
}

func RawJson(key string, value []byte) Field {
	return Field{Key: key, Type: RawJsonType, Value: value}
}

func Group(key string, fields ...Field) Field {
	return Field{Key: key, Type: GroupType, Value: fields}
}
