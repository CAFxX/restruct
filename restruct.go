package restruct

import (
	"encoding"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

type RawMatch string

func (r RawMatch) String() string {
	return string(r)
}

type group struct {
	index      int
	fieldIndex []int
}

func Compile[T any](restr string) func(s string) (T, error) {
	re := regexp.MustCompile(restr)

	var zero T
	t := reflect.TypeOf(zero)
	if t.Kind() != reflect.Struct {
		panic("type %s is not a struct")
	}

	tRawMatch := reflect.TypeOf(RawMatch(""))

	var fm []group

	raw, ok := t.FieldByName("RawMatch")
	if ok && raw.Type == tRawMatch {
		fm = append(fm, group{index: 0, fieldIndex: raw.Index})
	}

	for gi, gn := range re.SubexpNames() {
		f, ok := t.FieldByName(gn)
		if !ok {
			continue
		}
		switch f.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
			reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		case reflect.String:
			if f.Type == tRawMatch {
				panic("field " + gn + " is not a valid type")
			}
		default:
			if !f.Type.Implements(reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()) && !reflect.PointerTo(f.Type).Implements(reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()) {
				panic("field " + gn + " is not a valid type")
			}
		}
		fm = append(fm, group{index: gi, fieldIndex: f.Index})
	}

	return func(s string) (r T, err error) {
		m := re.FindStringSubmatch(s)
		if m == nil {
			var zero T
			return zero, errNoMatch{}
		}
		v := reflect.ValueOf(&r).Elem()
		for _, g := range fm {
			f := v.FieldByIndex(g.fieldIndex)
			err := unmarshalAndSet(f, m[g.index])
			if err != nil {
				var zero T
				return zero, fmt.Errorf("parsing field %s: %w", re.SubexpNames()[g.index], err)
			}
		}
		return r, nil
	}
}

type errNoMatch struct{}

func (errNoMatch) Error() string { return "no match" }

type errUnsupportedKind struct{}

func (errUnsupportedKind) Error() string { return "unsupported kind" }

func unmarshalAndSet(v reflect.Value, s string) error {
	switch v.Kind() {
	case reflect.Bool:
		b, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		v.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(s, 10, v.Type().Bits())
		if err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		i, err := strconv.ParseUint(s, 10, v.Type().Bits())
		if err != nil {
			return err
		}
		v.SetUint(i)
	case reflect.Float32, reflect.Float64:
		i, err := strconv.ParseFloat(s, v.Type().Bits())
		if err != nil {
			return err
		}
		v.SetFloat(i)
	case reflect.Complex64, reflect.Complex128:
		i, err := strconv.ParseComplex(s, v.Type().Bits())
		if err != nil {
			return err
		}
		v.SetComplex(i)
	case reflect.String:
		v.SetString(s)
	default:
		if u, ok := v.Interface().(encoding.TextUnmarshaler); ok {
			return u.UnmarshalText([]byte(s))
		} else if u, ok := v.Addr().Interface().(encoding.TextUnmarshaler); ok {
			return u.UnmarshalText([]byte(s))
		}
		return errUnsupportedKind{}
	}
	return nil
}
