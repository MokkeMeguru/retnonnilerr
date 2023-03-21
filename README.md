# retnonnilerr

`retnonnilerr` is a static tool to prevent below unhandling error.

```golang
func f() error {
    x, err := fn()
    if err != nil {
        return nil // !!!
    }
    fmt.Printf("x is %v\n", x)
    return nil
}
```

## How to use

### From CLI

```
cd path/to/product
go install github.com/MokkeMeguru/retnonnilerr/cmd/retnonnilerr
retnonnilerr ./...
```

## From CI

see. [my custom linter settings](./.github/workflows/ci.yml)

try it using [act](https://github.com/nektos/act)

```
act --job reviewdog
```

## Ignore Lint?

If you want to ignore this linter at the line, you can comment `lint:ignore retnonnilerr`.

```golang
func f() error {
    x, err := fn()
    if err != nil {
        // lint:ignore retnonnilerr TODO fix
        return nil
    }
    fmt.Printf("x is %v\n", x)
    return nil
}
```
