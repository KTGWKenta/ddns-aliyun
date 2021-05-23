package defines

import (
	"github.com/kentalee/errors"
)

const (
	ErrKeyReason  = "reason"
	ErrKeyContent = "content"
	ErrKeyType    = "type"
	ErrKeyField   = "field"
)

var (
	ErrInvalidDataLength   = errors.New("15fd649100010001", "invalid length")
	ErrInvalidDataRange    = errors.New("15fd649100010002", "invalid range")
	ErrInvalidDataType     = errors.New("15fd649100010003", "invalid value type")
	ErrFieldNotfound       = errors.New("15fd649100010004", "field not found")
	ErrValueNotfound       = errors.New("15fd649100010005", "value not found")
	ErrUninitializedObject = errors.New("15fd649100010006", "uninitialized object ( nil pointer )")
)
