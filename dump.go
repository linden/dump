package dump

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// format as: "{padding}{name} <{type}>".
func format(b *strings.Builder, v reflect.Value, prefix string, depth int) {
	vt := v.Type()

	// add padding.
	b.WriteString(strings.Repeat("\t", depth))

	// add prefix if applicable.
	if prefix != "" {
		b.WriteByte(' ')
		b.WriteString(prefix)
		b.WriteByte(' ')
	}

	// add the type.
	b.WriteByte('<')
	b.WriteString(vt.String())
	b.WriteByte('>')
}

func walk(b *strings.Builder, v reflect.Value, prefix string, depth int, stack []reflect.Value) {
	// ensure we're not printing something recursively.
	for _, c := range stack {
		if c == v {
			// add padding.
			b.WriteString(strings.Repeat("\t", depth))

			// add pseudo type.
			b.WriteString("<circular value>\n")

			return
		}
	}

	// add the current value to the stack.
	stack = append(stack, v)

	// format the current value.
	format(b, v, prefix, depth)

	// dereference the value until we have a non-pointer.
	for v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	kind := v.Kind()

	switch kind {
	// handle any types with fields.
	// TODO: handle arrays.
	case reflect.Struct, reflect.Map, reflect.Slice:
		// nothing else will be written on this line, move to the next.
		b.WriteByte('\n')

		depth += 1

		switch kind {
		case reflect.Struct:
			for i := 0; i < v.NumField(); i++ {
				name := v.Type().Field(i).Name
				f := v.Field(i)

				walk(b, f, name, depth, stack)
			}

		case reflect.Map:
			iter := v.MapRange()

			for iter.Next() {
				k := iter.Key()
				v := iter.Value()

				walk(b, v, fmt.Sprintf("%s:", k), depth, stack)
			}

		case reflect.Slice:

			for i := 0; i < v.Len(); i++ {
				itm := v.Index(i)

				walk(b, itm, fmt.Sprintf("%d:", i), depth, stack)
			}
		}

	// handle any primitive types.
	// TODO: handle interfaces (?), complex, floats and funcs.
	default:
		b.WriteByte(':')
		b.WriteByte(' ')

		switch kind {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			b.WriteString(strconv.Itoa(int(v.Int())))

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			b.WriteString(strconv.Itoa(int(v.Uint())))

		case reflect.Bool:
			if v.Bool() {
				b.WriteString("true")
			} else {
				b.WriteString("false")
			}

		default:
			b.WriteString(v.String())
		}

		b.WriteByte('\n')
	}
}

func Dump(x any) {
	b := new(strings.Builder)

	walk(b, reflect.ValueOf(x), "", 0, []reflect.Value{})

	fmt.Print(b.String())
}
