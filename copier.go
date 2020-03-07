package copier

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"reflect"
	"strings"
)

// These flags define options for tag handling
const (
	// Denotes that a destination field must be copied to. If copying fails then a panic will ensue.
	tagMust uint8 = 1 << iota

	// Denotes that the program should not panic when the must flag is on and
	// value is not copied. The program will return an error instead.
	tagNoPanic

	// Ignore a destation field from being copied to.
	tagIgnore

	// Denotes that the value as been copied
	hasCopied
)

// Option sets copy options
type Option struct {
	IgnoreEmpty bool
	DeepCopy    bool
}

// Copy copy things
func Copy(toValue interface{}, fromValue interface{}) (err error) {
	return copier(toValue, fromValue, Option{})
}

// CopyWithOption copy with option
func CopyWithOption(toValue interface{}, fromValue interface{}, opt Option) (err error) {
	return copier(toValue, fromValue, opt)
}

func copier(toValue interface{}, fromValue interface{}, opt Option) (err error) {
	var (
		isSlice bool
		amount  = 1
		from    = indirect(reflect.ValueOf(fromValue))
		to      = indirect(reflect.ValueOf(toValue))
	)

	if !to.CanAddr() {
		return ErrInvalidCopyDestination
	}

	// Return is from value is invalid
	if !from.IsValid() {
		return ErrInvalidCopyFrom
	}

	fromType, isPtrFrom := indirectType(from.Type())
	toType, _ := indirectType(to.Type())

	if fromType.Kind() == reflect.Interface {
		fromType = reflect.TypeOf(from.Interface())
	}

	if toType.Kind() == reflect.Interface {
		toType = reflect.TypeOf(to.Interface())
	}

	// Just set it if possible to assign for normal types
	if from.Kind() != reflect.Slice && from.Kind() != reflect.Struct && from.Kind() != reflect.Map && (from.Type().AssignableTo(to.Type()) || from.Type().ConvertibleTo(to.Type())) {
		if !isPtrFrom || !opt.DeepCopy {
			to.Set(from.Convert(to.Type()))
		} else {
			fromCopy := reflect.New(from.Type())
			fromCopy.Set(from.Elem())
			to.Set(fromCopy.Convert(to.Type()))
		}
		return
	}

	if fromType.Kind() == reflect.Map && toType.Kind() == reflect.Map {
		if !fromType.Key().ConvertibleTo(toType.Key()) {
			return ErrMapKeyNotMatch
		}

		if to.IsNil() {
			to.Set(reflect.MakeMapWithSize(toType, from.Len()))
		}

		for _, k := range from.MapKeys() {
			toKey := indirect(reflect.New(toType.Key()))
			if !set(toKey, k, opt.DeepCopy) {
				return fmt.Errorf("%w map, old key: %v, new key: %v", ErrNotSupported, k.Type(), toType.Key())
			}

			elemType, _ := indirectType(toType.Elem())
			toValue := indirect(reflect.New(elemType))
			if !set(toValue, from.MapIndex(k), opt.DeepCopy) {
				if err = copier(toValue.Addr().Interface(), from.MapIndex(k).Interface(), opt); err != nil {
					return err
				}
			}

			for {
				if elemType == toType.Elem() {
					to.SetMapIndex(toKey, toValue)
					break
				}
				elemType = reflect.PtrTo(elemType)
				toValue = toValue.Addr()
			}
		}
		return
	}

	if from.Kind() == reflect.Slice && to.Kind() == reflect.Slice && fromType.ConvertibleTo(toType) {
		if to.IsNil() {
			slice := reflect.MakeSlice(reflect.SliceOf(to.Type().Elem()), from.Len(), from.Cap())
			to.Set(slice)
		}
		for i := 0; i < from.Len(); i++ {
			if !set(to.Index(i), from.Index(i), opt.DeepCopy) {
				err = CopyWithOption(to.Index(i).Addr().Interface(), from.Index(i).Interface(), opt)
				if err != nil {
					continue
				}
			}
		}
		return
	}

	if fromType.Kind() != reflect.Struct || toType.Kind() != reflect.Struct {
		// skip not supported type
		return
	}

	if to.Kind() == reflect.Slice {
		isSlice = true
		if from.Kind() == reflect.Slice {
			amount = from.Len()
		}
	}

	for i := 0; i < amount; i++ {
		var dest, source reflect.Value

		if isSlice {
			// source
			if from.Kind() == reflect.Slice {
				source = indirect(from.Index(i))
			} else {
				source = indirect(from)
			}
			// dest
			dest = indirect(reflect.New(toType).Elem())
		} else {
			source = indirect(from)
			dest = indirect(to)
		}

		destKind := dest.Kind()
		initDest := false
		if destKind == reflect.Interface {
			initDest = true
			dest = indirect(reflect.New(toType))
		}

		// Get tag options
		tagBitFlags := map[string]uint8{}
		if dest.IsValid() {
			tagBitFlags = getBitFlags(toType)
		}

		// check source
		if source.IsValid() {
			// Copy from source field to dest field or method
			fromTypeFields := deepFields(fromType)
			for _, field := range fromTypeFields {
				name := field.Name

				// Get bit flags for field
				fieldFlags, _ := tagBitFlags[name]

				// Check if we should ignore copying
				if (fieldFlags & tagIgnore) != 0 {
					continue
				}

				if fromField := source.FieldByName(name); fromField.IsValid() && !shouldIgnore(fromField, opt.IgnoreEmpty) {
					// process for nested anonymous field
					destFieldNotSet := false
					if f, ok := dest.Type().FieldByName(name); ok {
						for idx, x := range f.Index {
							destFieldKind := dest.Field(x).Kind()
							if destFieldKind != reflect.Ptr {
								continue
							}

							if !dest.Field(x).IsNil() {
								continue
							}

							if !dest.Field(x).CanSet() {
								destFieldNotSet = true
								break
							}

							newValue := reflect.New(dest.FieldByIndex(f.Index[0 : idx+1]).Type().Elem())
							dest.Field(x).Set(newValue)
						}
					}

					if destFieldNotSet {
						break
					}

					toField := dest.FieldByName(name)
					if toField.IsValid() {
						if toField.CanSet() {
							if !set(toField, fromField, opt.DeepCopy) {
								if err := copier(toField.Addr().Interface(), fromField.Interface(), opt); err != nil {
									return err
								}
							} else {
								if fieldFlags != 0 {
									// Note that a copy was made
									tagBitFlags[name] = fieldFlags | hasCopied
								}
							}
						}
					} else {
						// try to set to method
						var toMethod reflect.Value
						if dest.CanAddr() {
							toMethod = dest.Addr().MethodByName(name)
						} else {
							toMethod = dest.MethodByName(name)
						}

						if toMethod.IsValid() && toMethod.Type().NumIn() == 1 && fromField.Type().AssignableTo(toMethod.Type().In(0)) {
							toMethod.Call([]reflect.Value{fromField})
						}
					}
				}
			}

			// Copy from from method to dest field
			for _, field := range deepFields(toType) {
				name := field.Name

				var fromMethod reflect.Value
				if source.CanAddr() {
					fromMethod = source.Addr().MethodByName(name)
				} else {
					fromMethod = source.MethodByName(name)
				}

				if fromMethod.IsValid() && fromMethod.Type().NumIn() == 0 && fromMethod.Type().NumOut() == 1 && !shouldIgnore(fromMethod, opt.IgnoreEmpty) {
					if toField := dest.FieldByName(name); toField.IsValid() && toField.CanSet() {
						values := fromMethod.Call([]reflect.Value{})
						if len(values) >= 1 {
							set(toField, values[0], opt.DeepCopy)
						}
					}
				}
			}
		}

		if isSlice {
			if dest.Addr().Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest.Addr()))
			} else if dest.Type().AssignableTo(to.Type().Elem()) {
				to.Set(reflect.Append(to, dest))
			}
		} else if initDest {
			to.Set(dest)
		}

		err = checkBitFlags(tagBitFlags)
	}

	return
}

func shouldIgnore(v reflect.Value, ignoreEmpty bool) bool {
	if !ignoreEmpty {
		return false
	}

	return v.IsZero()
}

func deepFields(reflectType reflect.Type) []reflect.StructField {
	if reflectType, _ = indirectType(reflectType); reflectType.Kind() == reflect.Struct {
		fields := make([]reflect.StructField, 0, reflectType.NumField())

		for i := 0; i < reflectType.NumField(); i++ {
			v := reflectType.Field(i)
			if v.Anonymous {
				fields = append(fields, deepFields(v.Type)...)
			} else {
				fields = append(fields, v)
			}
		}

		return fields
	}

	return nil
}

func indirect(reflectValue reflect.Value) reflect.Value {
	for reflectValue.Kind() == reflect.Ptr {
		reflectValue = reflectValue.Elem()
	}
	return reflectValue
}

func indirectType(reflectType reflect.Type) (_ reflect.Type, isPtr bool) {
	for reflectType.Kind() == reflect.Ptr || reflectType.Kind() == reflect.Slice {
		reflectType = reflectType.Elem()
		isPtr = true
	}
	return reflectType, isPtr
}

func set(to, from reflect.Value, deepCopy bool) bool {
	if from.IsValid() {
		if to.Kind() == reflect.Ptr {
			// set `to` to nil if from is nil
			if from.Kind() == reflect.Ptr && from.IsNil() {
				to.Set(reflect.Zero(to.Type()))
				return true
			} else if to.IsNil() {
				// `from`         -> `to`
				// sql.NullString -> *string
				if fromValuer, ok := driverValuer(from); ok {
					v, err := fromValuer.Value()
					if err != nil {
						return false
					}
					// if `from` is not valid do nothing with `to`
					if v == nil {
						return true
					}
				}
				// allocate new `to` variable with default value (eg. *string -> new(string))
				to.Set(reflect.New(to.Type().Elem()))
			}
			// depointer `to`
			to = to.Elem()
		}

		if deepCopy {
			toKind := to.Kind()
			if toKind == reflect.Interface && to.IsNil() {
				to.Set(reflect.New(reflect.TypeOf(from.Interface())).Elem())
				toKind = reflect.TypeOf(to.Interface()).Kind()
			}
			if toKind == reflect.Struct || toKind == reflect.Map || toKind == reflect.Slice {
				return false
			}
		}

		if from.Type().ConvertibleTo(to.Type()) {
			to.Set(from.Convert(to.Type()))
		} else if toScanner, ok := to.Addr().Interface().(sql.Scanner); ok {
			// `from`  -> `to`
			// *string -> sql.NullString
			if from.Kind() == reflect.Ptr {
				// if `from` is nil do nothing with `to`
				if from.IsNil() {
					return true
				}
				// depointer `from`
				from = indirect(from)
			}
			// `from` -> `to`
			// string -> sql.NullString
			// set `to` by invoking method Scan(`from`)
			err := toScanner.Scan(from.Interface())
			if err != nil {
				return false
			}
		} else if fromValuer, ok := driverValuer(from); ok {
			// `from`         -> `to`
			// sql.NullString -> string
			v, err := fromValuer.Value()
			if err != nil {
				return false
			}
			// if `from` is not valid do nothing with `to`
			if v == nil {
				return true
			}
			rv := reflect.ValueOf(v)
			if rv.Type().AssignableTo(to.Type()) {
				to.Set(rv)
			}
		} else if from.Kind() == reflect.Ptr {
			return set(to, from.Elem(), deepCopy)
		} else {
			return false
		}
	}

	return true
}

// parseTags Parses struct tags and returns uint8 bit flags.
func parseTags(tag string) (flags uint8) {
	for _, t := range strings.Split(tag, ",") {
		switch t {
		case "-":
			flags = tagIgnore
			return
		case "must":
			flags = flags | tagMust
		case "nopanic":
			flags = flags | tagNoPanic
		}
	}
	return
}

// getBitFlags Parses struct tags for bit flags.
func getBitFlags(toType reflect.Type) map[string]uint8 {
	flags := map[string]uint8{}
	toTypeFields := deepFields(toType)

	// Get a list dest of tags
	for _, field := range toTypeFields {
		tags := field.Tag.Get("copier")
		if tags != "" {
			flags[field.Name] = parseTags(tags)
		}
	}
	return flags
}

// checkBitFlags Checks flags for error or panic conditions.
func checkBitFlags(flagsList map[string]uint8) (err error) {
	// Check flag conditions were met
	for name, flags := range flagsList {
		if flags&hasCopied == 0 {
			switch {
			case flags&tagMust != 0 && flags&tagNoPanic != 0:
				err = fmt.Errorf("Field %s has must tag but was not copied", name)
				return
			case flags&(tagMust) != 0:
				panic(fmt.Sprintf("Field %s has must tag but was not copied", name))
			}
		}
	}
	return
}

func driverValuer(v reflect.Value) (i driver.Valuer, ok bool) {

	if !v.CanAddr() {
		i, ok = v.Interface().(driver.Valuer)
		return
	}

	i, ok = v.Addr().Interface().(driver.Valuer)
	return
}
