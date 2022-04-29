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

type Position int

type group struct {
	index      int
	fieldIndex []int
}

func Compile[T any](restr string) func(s string) (T, error) {
	re := regexp.MustCompile(restr)

	t := reflect.TypeOf(zero[T]())
	if t.Kind() != reflect.Struct {
		panic("type %s is not a struct")
	}

	tRawMatch := reflect.TypeOf(RawMatch(""))
	tPosition := reflect.TypeOf(Position(0))

	var fm []group

	if raw, ok := t.FieldByName("RawMatch"); ok && raw.Type == tRawMatch {
		fm = append(fm, group{index: 0, fieldIndex: raw.Index})
	}
	if pos, ok := t.FieldByName("Position"); ok && pos.Type == tPosition {
		fm = append(fm, group{index: -1, fieldIndex: pos.Index})
	}

	for gi, gn := range re.SubexpNames() {
		f, ok := t.FieldByName(gn)
		if !ok {
			continue
		}
		switch f.Type.Kind() {
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
			reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128:
		case reflect.Int:
			if f.Type == tPosition {
				panic("field " + gn + " is not a valid type")
			}
		case reflect.String:
			if f.Type == tRawMatch {
				panic("field " + gn + " is not a valid type")
			}
		case reflect.Slice, reflect.Array:
			if f.Type.Elem().Kind() != reflect.Uint8 {
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
		m := re.FindStringSubmatchIndex(s)
		if m == nil {
			return zero[T](), errNoMatch{}
		}
		v := reflect.ValueOf(&r).Elem()
		for _, g := range fm {
			f := v.FieldByIndex(g.fieldIndex)
			if g.index == -1 {
				f.SetInt(int64(m[0]))
				continue
			}
			err := unmarshalAndSet(f, s[m[g.index*2]:m[g.index*2+1]])
			if err != nil {
				return zero[T](), fmt.Errorf("parsing field %s: %w", re.SubexpNames()[g.index], err)
			}
		}
		return r, nil
	}
}

type errNoMatch struct{}

func (errNoMatch) Error() string { return "no match" }

type errUnsupportedKind struct{}

func (errUnsupportedKind) Error() string { return "unsupported kind" }

type errArrayOverflow struct{}

func (errArrayOverflow) Error() string { return "array too short" }

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
	case reflect.Slice:
		v.Set(reflect.ValueOf([]byte(s)))
	case reflect.Array:
		n := copy(v.Slice(0, v.Len()).Interface().([]byte), s)
		if n < len(s) {
			return errArrayOverflow{}
		}
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

func zero[T any]() (zero T) { return }
