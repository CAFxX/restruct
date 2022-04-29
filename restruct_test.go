package restruct

import (
	"fmt"
	"testing"
	"time"
)

func Example() {
	parse := Compile[struct {
		RawMatch
		Position
		Foo int
		Bar time.Time
	}](`(?P<Foo>[0-9]+)/(?P<Bar>[^\s]+)`)

	foo, err := parse("---> 42/2022-02-22T09:00:00Z <---")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(foo.Foo)
	fmt.Println(foo.Bar)
	fmt.Println(foo)
	fmt.Println(foo.Position)

	// Output:
	// 42
	// 2022-02-22 09:00:00 +0000 UTC
	// 42/2022-02-22T09:00:00Z
	// 5
}

func TestArray(t *testing.T) {
	r, err := Compile[struct {
		X [8]byte
	}](`(?P<X>[0-9]+)`)(`> 0123 <`)
	if err != nil {
		t.Fatal(err)
	}
	if r.X != [8]byte{'0', '1', '2', '3', 0, 0, 0, 0} {
		t.Fatal(r.X)
	}

	r, err = Compile[struct {
		X [8]byte
	}](`(?P<X>[0-9]+)`)(`> 012345678 <`)
	if err == nil {
		t.Fatal("nil error")
	}
}

func TestSlice(t *testing.T) {
	r, err := Compile[struct {
		X []byte
	}](`(?P<X>[0-9]+)`)(`> 0123 <`)
	if err != nil {
		t.Fatal(err)
	}
	if string(r.X) != "0123" {
		t.Fatal(r.X)
	}
}

func TestInt(t *testing.T) {
	r, err := Compile[struct {
		X int8
	}](`(?P<X>[0-9]+)`)(`> 0123 <`)
	if err != nil {
		t.Fatal(err)
	}
	if r.X != 123 {
		t.Fatal(r.X)
	}

	r, err = Compile[struct {
		X int8
	}](`(?P<X>[0-9]+)`)(`> 012345678 <`)
	if err == nil {
		t.Fatal("nil error")
	}
}

func TestUint(t *testing.T) {
	r, err := Compile[struct {
		X uint8
	}](`(?P<X>[0-9]+)`)(`> 0123 <`)
	if err != nil {
		t.Fatal(err)
	}
	if r.X != 123 {
		t.Fatal(r.X)
	}

	r, err = Compile[struct {
		X uint8
	}](`(?P<X>[0-9]+)`)(`> 012345678 <`)
	if err == nil {
		t.Fatal("nil error")
	}
}

func TestBool(t *testing.T) {
	r, err := Compile[struct {
		X bool
	}](`(?P<X>[0-9]+)`)(`> 1 <`)
	if err != nil {
		t.Fatal(err)
	}
	if r.X != true {
		t.Fatal(r.X)
	}

	r, err = Compile[struct {
		X bool
	}](`(?P<X>[0-9]+)`)(`> 0 <`)
	if err != nil {
		t.Fatal(err)
	}
	if r.X != false {
		t.Fatal(r.X)
	}

	r, err = Compile[struct {
		X bool
	}](`(?P<X>[0-9]+)`)(`> 5 <`)
	if err == nil {
		t.Fatal("nil error")
	}
}

func TestString(t *testing.T) {
	r, err := Compile[struct {
		X string
	}](`(?P<X>[0-9]+)`)(`> 0123 <`)
	if err != nil {
		t.Fatal(err)
	}
	if r.X != "0123" {
		t.Fatal(r.X)
	}
}

func TestTime(t *testing.T) {
	r, err := Compile[struct {
		X time.Time
	}](`(?P<X>[0-9TZ:.+-]+)`)(`> 2022-02-22T09:00:00.123+01:00 <`)
	if err != nil {
		t.Fatal(err)
	}
	exp, err := time.Parse(time.RFC3339Nano, `2022-02-22T09:00:00.123+01:00`)
	if err != nil {
		t.Fatal(err)
	}
	if !r.X.Equal(exp) {
		t.Fatal(r.X)
	}

	r, err = Compile[struct {
		X time.Time
	}](`(?P<X>[0-9]+)`)(`> 012345678 <`)
	if err == nil {
		t.Fatal("nil error")
	}
}

func TestNoMatch(t *testing.T) {
	_, err := Compile[struct {
		X string
	}](`(?P<X>[0-9]+)`)(`> xxx <`)
	exp := errNoMatch{}
	if err != exp {
		t.Fatal(err)
	}
}
