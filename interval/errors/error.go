package errors

import "errors"

var (
	maxLength       = 255
	ErrNilKey       = errors.New("key is nil")
	ErrKeyTooLong   = errors.New("key length is too long")
	ErrValueTooLong = errors.New("value length is too long")
)

var ErrReplayLog = errors.New("replay log err, idx out of range")

func CheckKey(k string) error {
	if k == "" {
		return ErrNilKey
	}
	if len(k) > maxLength {
		return ErrKeyTooLong
	}
	return nil
}

func CheckValue(v string) error {
	if len(v) > maxLength {
		return ErrValueTooLong
	}
	return nil
}
