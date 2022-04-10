package serialization

import "errors"

var (
	ErrSerializerAlreadyExists        = errors.New("register failed: serializer with name alredy exists")
	ErrSliceFiledSizeTagNotFound      = errors.New("(de)serialize failed: met struct who dont have a customized deserializer but contains slice that dont have a size field")
	ErrSerializeTypeDismatch          = errors.New("(de)serialize failed: type dismatch")
	ErrUsingUnaddressableValue        = errors.New("(de)serialize failed: using unaddressable value")
	ErrNonLastFieldContainsOmitOption = errors.New("(de)serialize failed: met struct contains a field with omit option but not last field of the struct")
)
