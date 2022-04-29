# restruct
Parse regular expressions to structs.

## Supported conversions

Matches can be parsed into the following field types:

- `string`
- `[]byte`
- `[N]byte`
- `int`/`uint`/`intN`/`uintN`/`uintptr`
- `floatN`
- `complexN`
- `bool`
- any type that implements [`encoding.TextUnmarshaler`](https://pkg.go.dev/encoding#TextUnmarshaler)

## Example

```golang
type Foo struct {
    Bar string
    Baz float64
}

parse, _ := Compile[Foo](`(?<Bar>[a-z]+) (?<Baz>[0-9.]+)`)

foo, _ := parse(`yadda 3.14`)

fmt.Printf("%s/%.3f", foo.Bar, foo.Baz) // Output: yadda/3.140
```