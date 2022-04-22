package copier

import "errors"

var (
	ErrInvalidCopyDestination        = errors.New("copy destination is invalid")
	ErrInvalidCopyFrom               = errors.New("copy from is invalid")
	ErrMapKeyNotMatch                = errors.New("map's key type doesn't match")
	ErrNotSupported                  = errors.New("not supported")
	ErrFieldNameTagStartNotUpperCase = errors.New("copier field name tag must be start upper case")
)
