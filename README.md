# retnonnilerr

[![Go Reference](https://pkg.go.dev/badge/github.com/MokkeMeguru/retnonnilerr.svg)](https://pkg.go.dev/github.com/MokkeMeguru/retnonnilerr)
[![Test](https://github.com/MokkeMeguru/retnonnilerr/actions/workflows/ci.yml/badge.svg)](https://github.com/MokkeMeguru/retnonnilerr/actions/workflows/ci.yml)

`retnonnilerr` is a static analysis tool to prevent below unhandling error.

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
go install github.com/MokkeMeguru/retnonnilerr/cmd/retnonnilerr
cd path/to/product
retnonnilerr ./...
```

## From CI

See [my custom linter settings](./.github/workflows/ci.yml)

Test it using [act](https://github.com/nektos/act) (WARNING: reviewdog cannot do normaly)

```
act --job reviewdog
```

## Ignore Lint?

If you want to ignore this linter at the line, you can comment `lint:ignore retnonnilerr`.

```golang
func f() error {
    x, err := fn()
    if err != nil {
        //lint:ignore retnonnilerr TODO fix
        return nil
    }
    fmt.Printf("x is %v\n", x)
    return nil
}
```

## References

- [nilerr](https://github.com/gostaticanalysis/nilerr/)

  Very similar, except for the lack of test code checks and the part that inspects err without exceptions.
