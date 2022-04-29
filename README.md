# restruct
Parse regular expressions to structs

## Example

```golang
type Foo struct {
    Bar string
    Baz float64
}

parse, _ := Compile[Foo](`(?<Bar>[a-z]+) (?<Baz>[0-9.]+)`)

foo, _ := parse(`yadda 3.14`)

fmt.Printf("%s/%f", foo.Bar, foo.Baz) // Output: yadda/3.14
```