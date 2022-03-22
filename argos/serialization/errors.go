package serialization

import "errors"

var (
	SerializerAlreadyExistsError   = errors.New("register failed: serializer with name alredy exists")
	SliceFiledSizeTagNotFound      = errors.New("(de)serialize failed: met struct who dont have a customized deserializer but contains slice that dont have a size field")
	SerializeTypeDismatchError     = errors.New("(de)serialize failed: type dismatch")
	UsingUnaddressableValue        = errors.New("(de)serialize failed: using unaddressable value")
	NonLastFieldContainsOmitOption = errors.New("(de)serialize failed: met struct contains a field with omit option but not last field of the struct")
)
