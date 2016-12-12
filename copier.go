package copier

import (
	"fmt"

	"reflect"
)

func Copy(toValue interface{}, fromValue interface{}) (err error) {
	var (
		isSlice   bool
		fromType  reflect.Type
		isFromPtr bool
		toType    reflect.Type
		amount    int
	)
	var accumulatedError error

	from := reflect.Indirect(reflect.ValueOf(fromValue))
	to := reflect.Indirect(reflect.ValueOf(toValue))

	if to.Kind() == reflect.Slice {
		isSlice = true
		if from.Kind() == reflect.Slice {
			fromType = from.Type().Elem()
			if fromType.Kind() == reflect.Ptr {
				fromType = fromType.Elem()
				isFromPtr = true
			}
			amount = from.Len()
		} else {
			fromType = from.Type()
			amount = 1
		}

		toType = to.Type().Elem()
		if toType.Kind() == reflect.Ptr {
			toType = toType.Elem()
		}
	} else {
		fromType = from.Type()
		toType = to.Type()
		amount = 1
	}

	if isSlice {
		if to.IsNil() {
			to.Set(reflect.MakeSlice(to.Type(), 0, amount))
		}

		if from.Kind() == reflect.Slice {
			if from.Type().Elem().Kind() == reflect.Ptr {
				newSlice := reflect.MakeSlice(to.Type(), amount, amount)
				originalLen := to.Len()
				to.Set(reflect.AppendSlice(to, newSlice))

				for i := 0; i < amount; i++ {
					var newT reflect.Value
					if to.Type().Elem().Kind() == reflect.Ptr {
						newT = reflect.New(to.Type().Elem().Elem())
					} else {
						newT = reflect.New(to.Type().Elem())
					}
					err := Copy(newT.Interface(), from.Index(i).Addr().Interface())
					to.Index(originalLen + i).Set(newT)
					if nil != err {
						if nil == accumulatedError {
							accumulatedError = err
							continue
						}
						accumulatedError = fmt.Errorf("error copying %v\n%v", err, accumulatedError)
					}
				}
			} else if from.Type().Elem().Kind() == reflect.Struct {
				newSlice := reflect.MakeSlice(to.Type(), amount, amount)
				originalLen := to.Len()
				to.Set(reflect.AppendSlice(to, newSlice))
				for i := 0; i < amount; i++ {
					err := Copy(to.Index(originalLen+i).Addr().Interface(), from.Index(i).Addr().Interface())
					if nil != err {
						if nil == accumulatedError {
							accumulatedError = err
							continue
						}
						accumulatedError = fmt.Errorf("error copying %v\n%v", err, accumulatedError)
					}
				}
			} else {
				reflect.Copy(to, from)
			}
		} else if from.Kind() == reflect.Struct {
			newSlice := reflect.MakeSlice(to.Type(), 1, 1)
			var newT reflect.Value
			if to.Type().Elem().Kind() == reflect.Ptr {
				newT = reflect.New(to.Type().Elem().Elem())
				newSlice.Index(0).Set(newT)
			} else {
				newT = reflect.New(to.Type().Elem())
				newSlice.Index(0).Set(newT.Elem())
			}
			originalLen := to.Len()
			to.Set(reflect.AppendSlice(to, newSlice))
			if to.Type().Elem().Kind() == reflect.Ptr {
				return Copy(to.Index(originalLen).Addr().Interface(), from.Addr().Interface())
			}

			return Copy(to.Index(originalLen).Addr().Interface(), from.Addr().Interface())
		} else if from.Kind() == reflect.Ptr {
			return Copy(toValue, from.Elem().Interface())
		}

		return fmt.Errorf("source slice type unsupported\n%v", accumulatedError)
	}

	for e := 0; e < amount; e++ {
		var dest, source reflect.Value
		if isSlice {
			if from.Kind() == reflect.Slice {
				source = from.Index(e)
				if isFromPtr {
					source = source.Elem()
				}
			} else {
				source = from
			}
		} else {
			source = from
		}

		if isSlice {
			dest = reflect.New(toType).Elem()
		} else {
			dest = to
		}

		for _, field := range deepFields(reflect.ValueOf(toValue).Type()) {
			name := field
			var fromField reflect.Value
			var fromMethod reflect.Value
			var toField reflect.Value
			var toMethod reflect.Value

			if source.Kind() == reflect.Ptr {
				if source.Elem().Kind() == reflect.Struct {
					fromField = source.Elem().FieldByName(name)
					fromMethod = source.MethodByName(name)
				} else {
					return fmt.Errorf("error\n%v", accumulatedError)
				}
			} else if source.Kind() == reflect.Struct {
				fromField = source.FieldByName(name)
				fromMethod = source.Addr().MethodByName(name)
			} else {
				return fmt.Errorf("error\n%v", accumulatedError)
			}

			if dest.Kind() == reflect.Ptr {
				if dest.Elem().Kind() == reflect.Struct {
					toField = dest.Elem().FieldByName(name)
					toMethod = dest.MethodByName(name)
				} else {
					return fmt.Errorf("error\n%v", accumulatedError)
				}
			} else if dest.Kind() == reflect.Struct {
				toField = dest.FieldByName(name)
				toMethod = dest.Addr().MethodByName(name)
			} else {
				return fmt.Errorf("error\n%v", accumulatedError)
			}

			canCopy := fromField.IsValid() && toMethod.IsValid() &&
				toMethod.Type().NumIn() == 1 && fromField.Type().AssignableTo(toMethod.Type().In(0))
			if canCopy {
				toMethod.Call([]reflect.Value{fromField})
				continue
			}

			canCopy = fromMethod.IsValid() && toField.IsValid() &&
				fromMethod.Type().NumOut() == 1 && fromMethod.Type().Out(0).AssignableTo(toField.Type())
			if canCopy {
				toField.Set(fromMethod.Call([]reflect.Value{})[0])
				continue
			}

			if fromMethod.IsValid() && toMethod.IsValid() {
			}
			canCopy = fromMethod.IsValid() && toMethod.IsValid() &&
				toMethod.Type().NumIn() == 1 && fromMethod.Type().NumOut() == 1 &&
				fromMethod.Type().Out(0).AssignableTo(toMethod.Type().In(0))
			if canCopy {
				toMethod.Call(fromMethod.Call([]reflect.Value{}))
				continue
			}

			_, accumulatedError = copyValue(toField, fromField, accumulatedError)
		}
	}
	return accumulatedError
}

func copyValue(to reflect.Value, from reflect.Value, accumulatedError error) (bool, error) {
	fieldsAreValid := to.IsValid() && from.IsValid()
	canCopy := fieldsAreValid && to.CanSet() && from.Type().AssignableTo(to.Type())

	if canCopy {
		to.Set(from)
		return true, accumulatedError
	}

	if !fieldsAreValid {
		return false, accumulatedError
	}

	_, accumulatedError = tryDeepCopyPtr(to, from, accumulatedError)
	_, accumulatedError = tryDeepCopyStruct(to, from, accumulatedError)
	_, accumulatedError = tryDeepCopySlice(to, from, accumulatedError)

	return false, accumulatedError
}

func tryDeepCopyPtr(toField reflect.Value, fromField reflect.Value, accumulatedError error) (bool, error) {
	deepCopyRequired := toField.Type().Kind() == reflect.Ptr && fromField.Type().Kind() == reflect.Ptr &&
		!fromField.IsNil() && toField.CanSet()

	copied := false
	if deepCopyRequired {
		toType := toField.Type().Elem()
		emptyObject := reflect.New(toType)
		toField.Set(emptyObject)
		err := Copy(toField.Interface(), fromField.Interface())
		if nil != err {
			copied = false
			if nil == accumulatedError {
				accumulatedError = err
				return false, accumulatedError
			}
			accumulatedError = fmt.Errorf("error copying %v\n%v", err, accumulatedError)
		} else {
			copied = true
		}
	}
	return copied, accumulatedError
}

func tryDeepCopyStruct(toField reflect.Value, fromField reflect.Value, accumulatedError error) (bool, error) {
	deepCopyRequired := toField.Type().Kind() == reflect.Struct && fromField.Type().Kind() == reflect.Struct && toField.CanSet()

	copied := false
	if deepCopyRequired {
		err := Copy(toField.Addr().Interface(), fromField.Addr().Interface())
		if nil != err {
			copied = false
			if nil == accumulatedError {
				accumulatedError = err
				return false, accumulatedError
			}
			accumulatedError = fmt.Errorf("error copying %v\n%v", err, accumulatedError)
		} else {
			copied = true
		}
	}
	return copied, accumulatedError
}

func tryDeepCopySlice(toField reflect.Value, fromField reflect.Value, accumulatedError error) (bool, error) {
	deepCopyRequired := toField.Type().Kind() == reflect.Slice && fromField.Type().Kind() == reflect.Slice && toField.CanSet()

	copied := false
	if deepCopyRequired {
		err := Copy(toField.Addr().Interface(), fromField.Addr().Interface())
		if nil != err {
			copied = false
			if nil == accumulatedError {
				accumulatedError = err
				return false, accumulatedError
			}
			accumulatedError = fmt.Errorf("error copying %v\n%v", err, accumulatedError)
		} else {
			copied = true
		}
	}
	return copied, accumulatedError
}

func deepFields(ifaceType reflect.Type) []string {
	fields := []string{}

	if ifaceType.Kind() == reflect.Ptr {
		// find all methods which take ptr as receiver
		fields = append(fields, deepFields(ifaceType.Elem())...)
	}

	// repeat (or do it for the first time) for all by-value-receiver methods
	fields = append(fields, deepFieldsImpl(ifaceType)...)

	return fields
}

func deepFieldsImpl(ifaceType reflect.Type) []string {
	fields := []string{}

	if ifaceType.Kind() != reflect.Ptr && ifaceType.Kind() != reflect.Struct {
		return fields
	}

	methods := ifaceType.NumMethod()
	for i := 0; i < methods; i++ {
		var v reflect.Method
		v = ifaceType.Method(i)

		fields = append(fields, v.Name)
	}

	if ifaceType.Kind() == reflect.Ptr {
		return fields
	}

	elements := ifaceType.NumField()
	for i := 0; i < elements; i++ {
		var v reflect.StructField
		v = ifaceType.Field(i)

		fields = append(fields, v.Name)
	}

	return fields
}
