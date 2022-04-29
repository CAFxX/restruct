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
		Baz []byte
	}](`(?P<Foo>[0-9]+)/(?P<Bar>[^\s]+)/(?P<Baz>x+)`)

	foo, err := parse("---> 42/2022-02-22T09:00:00Z/xxx <---")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(foo.Foo)
	fmt.Println(foo.Bar)
	fmt.Println(string(foo.Baz))
	fmt.Println(foo)
	fmt.Println(foo.Position)

	// Output:
	// 42
	// 2022-02-22 09:00:00 +0000 UTC
	// xxx
	// 42/2022-02-22T09:00:00Z/xxx
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
