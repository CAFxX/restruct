package restruct

import (
	"fmt"
	"time"
)

func Example() {
	parse := Compile[struct {
		RawMatch
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

	// Output:
	// 42
	// 2022-02-22 09:00:00 +0000 UTC
	// 42/2022-02-22T09:00:00Z
}
