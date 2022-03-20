package serialization

import "errors"

var (
	CannotCastToIntegerError       = errors.New("cast failed: cannot cast to integer")
	SliceFiledSizeTagNotFound      = errors.New("deserialize failed: met struct who dont have a customized deserializer but contains slice that dont have a size field")
	DeserializeTypeDismatchError   = errors.New("deserialize failed: type dismatch")
	NonLastFieldContainsOmitOption = errors.New("deserialize failed: met struct contains a field with omit option but not last field of the struct")
)
