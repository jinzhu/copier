package copier

import "reflect"

func Copy(copy_to interface{}, copy_from interface{}) (err error) {
	var (
		is_slice    bool
		from_typ    reflect.Type
		to_typ      reflect.Type
		elem_amount int
	)

	from := reflect.ValueOf(copy_from)
	to := reflect.ValueOf(copy_to)
	from_elem := reflect.Indirect(from)
	to_elem := reflect.Indirect(to)

	if to_elem.Kind() == reflect.Slice {
		is_slice = true
		if from_elem.Kind() == reflect.Slice {
			from_typ = from_elem.Type().Elem()
			elem_amount = from_elem.Len()
		} else {
			from_typ = from_elem.Type()
			elem_amount = 1
		}

		to_typ = to_elem.Type().Elem()
	} else {
		from_typ = from_elem.Type()
		to_typ = to_elem.Type()
		elem_amount = 1
	}

	for e := 0; e < elem_amount; e++ {
		var dest, source reflect.Value
		if is_slice {
			if from_elem.Kind() == reflect.Slice {
				source = from_elem.Index(e)
			} else {
				source = from_elem
			}
		} else {
			source = from_elem
		}

		if is_slice {
			dest = reflect.New(to_typ).Elem()
		} else {
			dest = to_elem
		}

		for i := 0; i < from_typ.NumField(); i++ {
			field := from_typ.Field(i)
			if !field.Anonymous {
				name := field.Name
				from_field := source.FieldByName(name)
				to_field := dest.FieldByName(name)
				to_method := dest.Addr().MethodByName(name)
				if from_field.IsValid() && to_field.IsValid() {
					to_field.Set(from_field)
				}

				if from_field.IsValid() && to_method.IsValid() {
					to_method.Call([]reflect.Value{from_field})
				}
			}
		}

		for i := 0; i < dest.NumField(); i++ {
			field := to_typ.Field(i)
			if !field.Anonymous {
				name := field.Name
				from_method := source.Addr().MethodByName(name)
				to_field := dest.FieldByName(name)

				if from_method.IsValid() && to_field.IsValid() {
					values := from_method.Call([]reflect.Value{})
					if len(values) >= 1 {
						to_field.Set(values[0])
					}
				}
			}
		}

		if is_slice {
			to_elem.Set(reflect.Append(to_elem, dest))
		}
	}
	return
}
