package form

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Values = url.Values

func Decode(values Values, fields Fields) error {
	var errs DecodeError
	for key, value := range fields {
		if err := value.DecodeForm(values[key]); err != nil {
			errs = append(errs, FieldError{Field: key, Err: err})
		}
	}
	if len(errs) > 0 {
		return errs
	}
	return nil
}

func Encode(fields Fields) Values {
	values := make(Values)
	for key, value := range fields {
		values[key] = value.EncodeForm()
	}
	return values
}

type FieldError struct {
	Field string `json:"field"`
	Err   error  `json:"error"`
}

func (e FieldError) Error() string {
	return fmt.Sprintf("failed to decode field %q: %v", e.Field, e.Err)
}

type DecodeError []FieldError

func (e DecodeError) Error() string {
	return "failed to decode form"
}

type Fields map[string]Value

type Value interface {
	DecodeForm([]string) error
	EncodeForm() []string
}

type GenericValue[T any] struct {
	Value      *T
	DecodeFunc func([]string) (T, bool, error)
	EncodeFunc func(T) []string
	EqualFunc  func(T, T) bool

	defaultVal *T
}

func (v GenericValue[T]) DecodeForm(values []string) error {
	if v.Value == nil {
		panic("form: nil value")
	}
	if v.DecodeFunc == nil {
		panic("form: nil decode func")
	}
	value, ok, err := v.DecodeFunc(values)
	if err != nil {
		return err
	}
	if ok {
		*v.Value = value
	} else if v.defaultVal != nil {
		*v.Value = *v.defaultVal
	}
	return nil
}

func (v GenericValue[T]) EncodeForm() []string {
	if v.Value == nil {
		panic("form: nil value")
	}
	if v.EncodeFunc == nil {
		panic("form: nil encode func")
	}
	if v.EqualFunc == nil {
		panic("form: nil equal func")
	}
	// don't encode zero or default values
	var zero T
	if v.EqualFunc(*v.Value, zero) {
		return nil
	}
	if v.defaultVal != nil && v.EqualFunc(*v.Value, *v.defaultVal) {
		return nil
	}
	return v.EncodeFunc(*v.Value)
}

func (v GenericValue[T]) Default(value T) GenericValue[T] {
	v.defaultVal = &value
	return v
}

func String[T ~string](ptr *T) GenericValue[T] {
	return GenericValue[T]{
		Value:      ptr,
		DecodeFunc: DecodeString[T],
		EncodeFunc: EncodeString[T],
		EqualFunc:  equal[T],
	}
}

func equal[T comparable](a, b T) bool {
	return a == b
}

func DecodeString[T ~string](values []string) (T, bool, error) {
	if len(values) == 0 {
		return "", false, nil
	}
	for _, strval := range values {
		strval = strings.TrimSpace(strval)
		if strval == "" {
			return "", true, nil
		}
		return T(strval), true, nil
	}
	return "", false, nil
}

func EncodeString[T ~string](value T) []string {
	return []string{string(value)}
}

type intConstraint interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

func Int[T intConstraint](ptr *T) GenericValue[T] {
	return GenericValue[T]{
		Value:      ptr,
		DecodeFunc: DecodeInt[T],
		EncodeFunc: EncodeInt[T],
		EqualFunc:  equal[T],
	}
}

func DecodeInt[T intConstraint](values []string) (T, bool, error) {
	if len(values) == 0 {
		return 0, false, nil
	}
	for _, strval := range values {
		strval = strings.TrimSpace(strval)
		if strval == "" {
			return 0, true, nil
		}
		value, err := strconv.ParseInt(strval, 10, 64)
		if err != nil {
			return 0, false, err
		}
		return T(value), true, nil
	}
	return 0, false, nil
}

func EncodeInt[T intConstraint](value T) []string {
	return []string{strconv.FormatInt(int64(value), 10)}
}

type uintConstraint interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func Uint[T uintConstraint](ptr *T) GenericValue[T] {
	return GenericValue[T]{
		Value:      ptr,
		DecodeFunc: DecodeUint[T],
		EncodeFunc: EncodeUint[T],
		EqualFunc:  equal[T],
	}
}

func DecodeUint[T uintConstraint](values []string) (T, bool, error) {
	if len(values) == 0 {
		return 0, false, nil
	}
	for _, strval := range values {
		strval = strings.TrimSpace(strval)
		if strval == "" {
			return 0, true, nil
		}
		value, err := strconv.ParseUint(strval, 10, 64)
		if err != nil {
			return 0, false, err
		}
		return T(value), true, nil
	}
	return 0, false, nil
}

func EncodeUint[T uintConstraint](value T) []string {
	return []string{strconv.FormatUint(uint64(value), 10)}
}

type floatConstraint interface {
	~float32 | ~float64
}

func Float[T floatConstraint](ptr *T) GenericValue[T] {
	return GenericValue[T]{
		Value:      ptr,
		DecodeFunc: DecodeFloat[T],
		EncodeFunc: EncodeFloat[T],
		EqualFunc:  equal[T],
	}
}

func DecodeFloat[T floatConstraint](values []string) (T, bool, error) {
	if len(values) == 0 {
		return 0, false, nil
	}
	for _, strval := range values {
		strval = strings.TrimSpace(strval)
		if strval == "" {
			return 0, true, nil
		}
		value, err := strconv.ParseFloat(strval, 64)
		if err != nil {
			return 0, false, err
		}
		return T(value), true, nil
	}
	return 0, false, nil
}

func EncodeFloat[T floatConstraint](value T) []string {
	return []string{strconv.FormatFloat(float64(value), 'f', -1, 64)}
}

func Time(ptr *time.Time, layout string) GenericValue[time.Time] {
	return GenericValue[time.Time]{
		Value:      ptr,
		DecodeFunc: DecodeTime(layout),
		EncodeFunc: EncodeTime(layout),
		EqualFunc:  time.Time.Equal,
	}
}

func DecodeTime(layout string) func([]string) (time.Time, bool, error) {
	return func(values []string) (time.Time, bool, error) {
		if len(values) == 0 {
			return time.Time{}, false, nil
		}
		for _, strval := range values {
			strval = strings.TrimSpace(strval)
			if strval == "" {
				return time.Time{}, true, nil
			}
			value, err := time.Parse(layout, strval)
			if err != nil {
				return time.Time{}, false, err
			}
			return value, true, nil
		}
		return time.Time{}, false, nil
	}
}

func EncodeTime(layout string) func(time.Time) []string {
	return func(value time.Time) []string {
		return []string{value.Format(layout)}
	}
}

func Bool[T ~bool](ptr *T) GenericValue[T] {
	return GenericValue[T]{
		Value:      ptr,
		DecodeFunc: DecodeBool[T],
		EncodeFunc: EncodeBool[T],
		EqualFunc:  equal[T],
	}
}

func DecodeBool[T ~bool](values []string) (T, bool, error) {
	if len(values) == 0 {
		return false, false, nil
	}
	for _, strval := range values {
		strval = strings.TrimSpace(strval)
		if strval == "" {
			return false, true, nil
		}
		value, err := strconv.ParseBool(strval)
		if err != nil {
			return false, false, err
		}
		return T(value), true, nil
	}
	return false, false, nil
}

func EncodeBool[T ~bool](value T) []string {
	return []string{strconv.FormatBool(bool(value))}
}
