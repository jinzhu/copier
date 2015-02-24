package copier

import "reflect"

func Copy(toValue interface{}, fromValue interface{}) (err error) {
	var (
		isSlice   bool
		fromType  reflect.Type
		isFromPtr bool
		toType    reflect.Type
		isToPtr   bool
		amount    int
	)

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
			isToPtr = true
		}
	} else {
		fromType = from.Type()
		toType = to.Type()
		amount = 1
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

		for i := 0; i < fromType.NumField(); i++ {
			field := fromType.Field(i)
			if !field.Anonymous {
				name := field.Name
				fromField := source.FieldByName(name)
				toField := dest.FieldByName(name)
				toMethod := dest.Addr().MethodByName(name)
				if fromField.IsValid() && toField.IsValid() {
					toField.Set(fromField)
				}

				if fromField.IsValid() && toMethod.IsValid() {
					toMethod.Call([]reflect.Value{fromField})
				}
			}
		}

		for i := 0; i < toType.NumField(); i++ {
			field := toType.Field(i)
			if !field.Anonymous {
				name := field.Name
				fromMethod := source.Addr().MethodByName(name)
				toField := dest.FieldByName(name)

				if fromMethod.IsValid() && toField.IsValid() {
					values := fromMethod.Call([]reflect.Value{})
					if len(values) >= 1 {
						toField.Set(values[0])
					}
				}
			}
		}

		if isSlice {
			if isToPtr {
				to.Set(reflect.Append(to, dest.Addr()))
			} else {
				to.Set(reflect.Append(to, dest))
			}
		}
	}
	return
}
