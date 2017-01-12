package copier

import (
	"fmt"

	"reflect"
	"strings"
)

// Copy copies one value to another. If you read this you probably came here for this function :-)
func Copy(toValue interface{}, fromValue interface{}) (err error) {
	from := reflect.ValueOf(fromValue)
	to := reflect.ValueOf(toValue)

	return copyImpl(to, from)
}

// copyImpl main function implementation.
// Differs only in accepting type to prevent pointer-to-value convertion each time
func copyImpl(to reflect.Value, from reflect.Value) error {
	var (
		isSlice  bool
		isStruct bool
	)
	var accumulatedError error

	toReducted := reductPointers(to)
	fromReducted := reductPointers(from)

	toKind := to.Kind()
	fromKind := from.Kind()

	toType := deepType(toReducted.Type())
	fromType := deepType(fromReducted.Type())

	toExact := exactValue(toReducted)
	fromExact := exactValue(fromReducted)

	toDepth := depth(to.Type())
	fromDepth := depth(from.Type())

	// different levels of indirection. not an error, though
	if toDepth != fromDepth {
		return nil
	}

	if toKind == reflect.Slice {
		isSlice = true
	}

	if toKind == reflect.Struct || toReducted.IsValid() && toReducted.Kind() == reflect.Ptr && toReducted.Elem().IsValid() && toReducted.Type().Elem().Kind() == reflect.Struct {
		isStruct = true
	}

	// destination is a slice
	if isSlice {
		// length without any indirection
		amount := deepLen(from)

		err := copySliceImpl(toExact, toKind, toType, fromExact, fromKind, fromType, amount, nil)
		if nil != err {
			return err
		}
		return nil
	}

	// destination is a struct
	if isStruct {
		for _, name := range deepFields(toReducted.Type()) {
			fromField := fieldByName(fromReducted, name)
			fromMethod := methodByName(fromReducted, name)

			var from reflect.Value

			if fromField.IsValid() {
				from = fromField
			} else if fromMethod.IsValid() && fromMethod.Type().NumOut() == 1 {
				from = fromMethod.Call([]reflect.Value{})[0]
			} else {
				continue
			}

			toField := fieldByName(toReducted, name)

			// if struct field is a slice we must create it here
			if toField.IsValid() && toField.Kind() == reflect.Slice && toField.IsNil() {
				capacity := 1
				if from.Kind() == reflect.Slice {
					capacity = from.Len()
				}
				// invoke the same method as for a root-level slice
				err := copySliceImpl(toField, toField.Kind(), toField.Type(), from, from.Kind(), from.Type(), capacity, nil)
				if nil != err {
					if nil == accumulatedError {
						accumulatedError = err
					} else {
						accumulatedError = fmt.Errorf("%v\n%v", err, accumulatedError)
					}
				}

				continue
			}

			toMethod := methodByName(toReducted, name)

			// we can't make stuff like deep copies when copying to a method
			canCopy := from.IsValid() && toMethod.IsValid() && toMethod.Kind() == reflect.Func && toMethod.Type().NumIn() == 1 && from.Type().AssignableTo(toMethod.Type().In(0))
			if canCopy {
				toMethod.Call([]reflect.Value{from})
				continue
			}

			_, accumulatedError = copyValue(toField, from, accumulatedError)
		}
		return accumulatedError
	}

	var err error
	_, err = copyValue(to, from, accumulatedError)
	return err
}

// fieldByName get field by name if possible
func fieldByName(base reflect.Value, name string) reflect.Value {
	if !base.IsValid() {
		return reflect.Zero(base.Type())
	}

	if base.Kind() == reflect.Ptr && !base.IsNil() {
		return fieldByName(base.Elem(), name)
	}
	if base.Kind() == reflect.Struct {
		return base.FieldByName(name)
	}
	return reflect.Zero(base.Type())
}

// methodByName get method by name if possible
func methodByName(base reflect.Value, name string) reflect.Value {
	if !base.IsValid() {
		return reflect.Zero(base.Type())
	}

	if base.Kind() == reflect.Ptr && !base.IsNil() {
		if base.Elem().Kind() == reflect.Struct {
			result := base.MethodByName(name)
			if result.IsValid() {
				return result
			}
		}
		return methodByName(base.Elem(), name)
	}
	if base.Kind() == reflect.Struct {
		return base.MethodByName(name)
	}
	return reflect.Zero(base.Type())
}

// copyValue tries to copy field using all supported methods
func copyValue(to reflect.Value, from reflect.Value, accumulatedError error) (bool, error) {
	if !to.IsValid() {
		return false, fmt.Errorf("destination is invalid")
	}

	if !from.IsValid() {
		return false, fmt.Errorf("source is invalid")
	}

	// this copy will work if and only if both values are primitive types
	err := tryCopyPrimitive(to, from)

	if err == nil {
		return true, nil
	}

	copied, accumulatedError := tryDeepCopyPtr(to, from, accumulatedError)
	if !copied {
		copied, accumulatedError = tryDeepCopyStruct(to, from, accumulatedError)
	}

	return copied, accumulatedError
}

// tryCopyPrimitive is invoked when direct assignment is far enough
func tryCopyPrimitive(dest reflect.Value, src reflect.Value) error {
	if !dest.IsValid() {
		return fmt.Errorf("destination is invalid")
	}

	if !src.IsValid() {
		return fmt.Errorf("source is invalid")
	}

	if !dest.CanSet() {
		return fmt.Errorf("destination is not settable")
	}

	if !src.Type().AssignableTo(dest.Type()) {
		return fmt.Errorf("destination type %v is not compatible with source type %v", dest.Type(), src.Type())
	}

	if dest.Kind() == reflect.Ptr {
		return fmt.Errorf("destination type %v is a pointer, pointers must not be assigned in this way", dest.Type())
	}

	dest.Set(src)

	return nil
}

// tryDeepCopyPtr goes through the chain of pointers, creating missing objects and, finally, invoking the recursion on the last step
func tryDeepCopyPtr(toField reflect.Value, fromField reflect.Value, accumulatedError error) (bool, error) {
	toDepth := depth(toField.Type())
	fromDepth := depth(fromField.Type())

	deepCopyRequired := toField.Type().Kind() == reflect.Ptr && fromField.Type().Kind() == reflect.Ptr &&
		!fromField.IsNil() && (toField.CanSet() || !toField.IsNil()) && toDepth == fromDepth && toField.IsValid() && fromField.IsValid()

	copied := false
	if deepCopyRequired {
		fromField = reductPointers(fromField)

		for toField.IsValid() && toField.Kind() == reflect.Ptr && toField.IsNil() {
			if !toField.CanSet() {
				return false, fmt.Errorf("cannot set empty pointer")
			}

			newTo := reflect.New(toField.Type().Elem())
			toField.Set(newTo)
			toField = reductPointers(toField)
		}

		toExact := exactValue(toField)
		fromExact := exactValue(fromField)

		var err error

		// structures may contain methods, which are connected with struct pointer, not by value
		if toExact.Kind() == reflect.Struct {
			err = copyImpl(toField, fromField)
		} else {
			err = copyImpl(toExact, fromExact)
		}

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

// tryDeepCopyStruct tries to copy field as a structure
func tryDeepCopyStruct(toField reflect.Value, fromField reflect.Value, accumulatedError error) (bool, error) {
	deepCopyRequired := toField.Type().Kind() == reflect.Struct && fromField.Type().Kind() == reflect.Struct && toField.CanSet()

	copied := false
	if toField.Kind() == reflect.Ptr && toField.IsNil() {
	}
	if deepCopyRequired {
		err := copyImpl(toField, fromField)
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

// copySliceImpl copies one slice to another, taking cate about empty slices and their types
func copySliceImpl(toExact reflect.Value, toKind reflect.Kind, toType reflect.Type,
	fromExact reflect.Value, fromKind reflect.Kind, fromType reflect.Type, amount int, accumulatedError error) error {

	if toKind != reflect.Slice {
		return fmt.Errorf("destination is not a slice: %v", toKind)
	}

	if toExact.IsNil() {
		toExact.Set(reflect.MakeSlice(toType, 0, 0))
	}

	originalLen := toExact.Len()

	newSlice := reflect.MakeSlice(toType, amount, amount)
	toExact.Set(reflect.AppendSlice(toExact, newSlice))

	if fromKind == reflect.Slice {
		for i := 0; i < amount; i++ {
			newT := reflect.New(toType.Elem())
			newIndirected := reflect.Indirect(newT)

			var destination reflect.Value
			if newIndirected.Kind() == reflect.Struct {
				destination = newT
			} else {
				destination = newIndirected
			}
			source := structAddr(fromExact.Index(i))

			err := copyImpl(destination, source)

			toExact.Index(originalLen + i).Set(newIndirected)

			if nil != err {
				if nil == accumulatedError {
					accumulatedError = err
					continue
				}
				accumulatedError = fmt.Errorf("error copying %v\n%v", err, accumulatedError)
			}
		}
	} else {
		newT := reflect.New(toType.Elem())
		newIndirected := reflect.Indirect(newT)

		var destination reflect.Value
		if newIndirected.Kind() == reflect.Struct {
			destination = newT
		} else {
			destination = newIndirected
		}
		source := structAddr(fromExact)

		err := copyImpl(destination, source)
		toExact.Index(originalLen).Set(newIndirected)

		if nil != err {
			if nil == accumulatedError {
				return err
			}
			return fmt.Errorf("error copying %v\n%v", err, accumulatedError)
		}

		return fmt.Errorf("source slice type unsupported\n%v", accumulatedError)
	}

	return nil
}

// structAddr returns addr of a structure value if it is possible to get it
func structAddr(val reflect.Value) reflect.Value {
	if val.IsValid() && val.Kind() == reflect.Struct && val.CanAddr() {
		return val.Addr()
	}
	return val
}

// deepFields returns all public field and methods names
func deepFields(ifaceType reflect.Type) []string {
	fields := []string{}

	if ifaceType.Kind() == reflect.Ptr && ifaceType.Elem().Kind() == reflect.Struct {
		// find all methods which take ptr as receiver
		fields = append(fields, deepFields(ifaceType.Elem())...)
	}

	// repeat (or do it for the first time) for all by-value-receiver methods
	fields = append(fields, deepFieldsImpl(ifaceType)...)

	return fields
}

func deepFieldsImpl(ifaceType reflect.Type) []string {
	fields := []string{}

	if ifaceType.Kind() != reflect.Ptr && ifaceType.Kind() != reflect.Struct ||
		ifaceType.Kind() == reflect.Ptr && ifaceType.Elem().Kind() == reflect.Slice {
		return fields
	}

	methods := ifaceType.NumMethod()
	for i := 0; i < methods; i++ {
		var v reflect.Method
		v = ifaceType.Method(i)

		if len(v.Name) == 0 || v.Name[0:1] != strings.ToUpper(v.Name[0:1]) {
			continue
		}

		fields = append(fields, v.Name)
	}

	if ifaceType.Kind() == reflect.Ptr {
		return fields
	}

	elements := ifaceType.NumField()
	for i := 0; i < elements; i++ {
		var v reflect.StructField
		v = ifaceType.Field(i)

		if len(v.Name) == 0 || v.Name[0:1] != strings.ToUpper(v.Name[0:1]) {
			continue
		}

		fields = append(fields, v.Name)
	}

	return fields
}

// deepLen returns amount of values after pointer chain.
// Amount is equal to slice length or to 1 if deepest type is not a slice
func deepLen(array reflect.Value) int {
	if array.IsValid() && array.Kind() == reflect.Slice {
		return array.Len()
	} else if array.IsValid() && array.Kind() == reflect.Ptr {
		return deepLen(array.Elem())
	}

	return 1
}

// deepKind returns kind of a deepest non-pointer type
func deepKind(ptr reflect.Type) reflect.Kind {
	if ptr.Kind() == reflect.Ptr {
		return deepKind(ptr.Elem())
	}

	return ptr.Kind()
}

// depth returns length of a pointer chain
func depth(ptr reflect.Type) int {
	if ptr.Kind() == reflect.Ptr {
		return 1 + depth(ptr.Elem())
	}

	return 0
}

// deepType returns type of first non-pointer
func deepType(ptrType reflect.Type) reflect.Type {
	if ptrType.Kind() == reflect.Ptr {
		return deepType(ptrType.Elem())
	}

	return ptrType
}

// reductPointers eliminates all pointers except the last one.
// Results of this function may be useful to obtain both structure fields and methods
func reductPointers(ptr reflect.Value) reflect.Value {
	if ptr.IsValid() && ptr.Kind() == reflect.Ptr && ptr.Elem().Kind() == reflect.Ptr {
		return reductPointers(ptr.Elem())
	}
	return ptr
}

func exactValue(ptr reflect.Value) reflect.Value {
	if ptr.Kind() == reflect.Ptr {
		return exactValue(ptr.Elem())
	}
	return ptr
}
