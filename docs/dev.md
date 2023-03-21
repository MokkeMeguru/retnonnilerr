## Project initialize

references:

- [【Go 言語】自作した静的解析ツールを GitHub Actions で実行する](https://zenn.dev/ytakaya/articles/55a07808c2fd5e)

```
.
├── README.md
├── cmd // cli で呼び出すツール
│   └── {package-name}
│       └── main.go
├── go.mod
├── {package-name}.go // 本体のコード
├── {package-name}_test.go
└── testdata // テストデータ
    └── src
        └── a
            ├── a.go
            └── a_test.go

```
