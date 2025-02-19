package emptycheck

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"
)

type RequiredFieldError struct {
	Field string
}

func (e RequiredFieldError) Error() string {
	return fmt.Sprintf("required field %s is empty", e.Field)
}

// 構造体のフィールドのすべてがemptyでないことをチェックする。
// タグでrequire:"noRequired"としている場合は、チェックをしない。
// フィールドはすべてExported fieldである必要がある。
func EmptyCheck[S any](s *S) error {
	rv := reflect.ValueOf(s).Elem()
	if rv.Kind() != reflect.Struct {
		panic("arg must be pointer to struct")
	}
	return emptyCheck(rv, "")
}

func emptyCheck(rv reflect.Value, parentTagName string) error {
	if rv.Kind() == reflect.Slice || rv.Kind() == reflect.Array {
		if rv.Len() == 0 {
			return RequiredFieldError{
				Field: parentTagName,
			}
		}
		if _, ok := rv.Interface().(uuid.UUID); ok {
			if isStructFieldEmpty(rv) {
				return RequiredFieldError{
					Field: parentTagName,
				}
			}
			return nil
		}
		for i := 0; i < rv.Len(); i++ {
			if rv.Index(i).Kind() == reflect.Struct {
				switch rv.Index(i).Interface().(type) {
				case Zeroable, Emptiable:
					if isStructFieldEmpty(rv.Field(i)) {
						return RequiredFieldError{
							Field: parentTagName,
						}
					}
				default:
					if err := emptyCheck(rv.Index(i), parentTagName); err != nil {
						return err
					}
				}
			} else if rv.Index(i).Kind() == reflect.Slice || rv.Index(i).Kind() == reflect.Array {
				if err := emptyCheck(rv.Index(i), parentTagName); err != nil {
					return err
				}
			}
		}
		return nil
	}

	if rv.Kind() == reflect.Struct {
		rt := rv.Type()
		for i := 0; i < rt.NumField(); i++ {
			//fmt.Println(rt.Field(i).Name)
			v := rt.Field(i).Tag.Get("require")
			if v == "noRequired" {
				continue
			}
			if rv.Field(i).Kind() == reflect.Struct {
				switch rv.Field(i).Interface().(type) {
				case Zeroable, Emptiable:
					if isStructFieldEmpty(rv.Field(i)) {
						return RequiredFieldError{
							Field: rt.Field(i).Name,
						}
					}
				default:
					err := emptyCheck(rv.Field(i), rt.Field(i).Name)
					if err != nil {
						return err
					}
				}
			} else if rv.Field(i).Kind() == reflect.Slice || rv.Field(i).Kind() == reflect.Array {
				err := emptyCheck(rv.Field(i), rt.Field(i).Name)
				if err != nil {
					return err
				}
			} else {
				if isStructFieldEmpty(rv.Field(i)) {
					return RequiredFieldError{
						Field: rt.Field(i).Name,
					}
				}
			}
		}
		return nil
	}

	panic(fmt.Sprintf("unexpected type: %T", rv.Interface()))
}

// time.Timeや*time.Time（nilではないもの）が該当する
type Zeroable interface {
	IsZero() bool
}

// Emptiableを実装することで独自にEmptyを定義可能
type Emptiable interface {
	IsEmpty() bool
}

// boolについては含めていない。
func isStructFieldEmpty(rv reflect.Value) bool {
	if rv.Kind() == reflect.Ptr {
		return rv.IsNil()
	}
	switch v := rv.Interface().(type) {
	case string:
		return v == ""
	case int:
		return v == 0
	case uint:
		return v == 0
	case int8:
		return v == 0
	case int16:
		return v == 0
	case int32:
		return v == 0
	case int64:
		return v == 0
	case float32:
		return v == 0
	case float64:
		return v == 0
	case uuid.UUID:
		return v == uuid.Nil
	case Zeroable:
		return v.IsZero()
	case Emptiable:
		return v.IsEmpty()
	default:
		panic(fmt.Sprintf("unexpected type: %T", v))
	}
}
