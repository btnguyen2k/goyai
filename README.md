# goyai

[![Go Report Card](https://goreportcard.com/badge/github.com/btnguyen2k/goyai)](https://goreportcard.com/report/github.com/btnguyen2k/goyai)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/btnguyen2k/goyai)](https://pkg.go.dev/github.com/btnguyen2k/goyai)
[![Actions Status](https://github.com/btnguyen2k/goyai/workflows/goyai/badge.svg)](https://github.com/btnguyen2k/goyai/actions)
[![codecov](https://codecov.io/gh/btnguyen2k/goyai/branch/main/graph/badge.svg?token=x12xW1YfiY)](https://codecov.io/gh/btnguyen2k/goyai)
[![Release](https://img.shields.io/github/release/btnguyen2k/goyai.svg?style=flat-square)](RELEASE-NOTES.md)

Yet Another I18n package for Go (Golang).

## Feature overview

- Support localized text messages, with plural forms.
- Support template string with named variables following [text/template](http://golang.org/pkg/text/template/) syntax.
- Support language files in JSON and YAML formats.

## Installation

```go
go get github.com/btnguyen2k/goyai
```

## Usage & Documentation

Package documentation is available at https://pkg.go.dev/github.com/btnguyen2k/goyai.

**Language file format**

`goyai` supports language file in JSON or YAML formats. Here is an example of language file in YAML:

```yaml
---
en:
  _name: English
  hello: Hello, world!
  remaining_tasks:
    zero: Congratulation! You are free now.
    one: There is 1 task left.
    two: There are 2 tasks left.
    few: There are a few tasks left.
    many: There are many tasks left.
    other: Hmmm!

vi:
  _name: Tiếng Việt
  hello: Xin chào
  remaining_tasks:
    zero: Chúc mừng! Bạn đã hoàn thành công việc.
    one: Vẫn còn 1 công việc nữa cần hoàn thành.
    two: Còn 1 công việc nữa cần hoàn thành.
    few: Vẫn còn vài công việc nữa cần hoàn thành.
    many: Còn quá trời việc chưa hoàn thành.
    other: Á chà!
```

> Multi-document YAML is currently **not** supported! Only the first document in multi-document YAML file is loaded.

**Load language files and build an I18n instance to use**

```go
import "github.com/btnguyen2k/goyai"

i18n, err := BuildI18n(goyai.I18nOptions{ConfigFileOrDir: "all_in_one.yaml", goyai.I18nFileFormat: goyai.Yaml, DefaultLocale: "en"})
if err != nil {
    panic(err)
}
```

> Supported language file formats are `goyai.Json` and `goyai.Yaml`. File format is automatically detected if `goyai.Auto` is specified.

Language messages can spread multiple files in a directory and be loaded in one go:

```go
i18n, err := BuildI18n(goyai.I18nOptions{ConfigFileOrDir: "./languages/", goyai.I18nFileFormat: goyai.Auto, DefaultLocale: "en"})
```

**Localize messages via I18n instance**

```go
// output "Xin chào"
fmt.Println(i18n.Localize("vi", "hello"))

// locale "nf" not exist, fall back to default locale "en"
// output "Hello, world!"
fmt.Println(i18n.Localize("nf", "hello"))

// output "Congratulation! You are free now."
fmt.Println(i18n.Localize("en", "remaining_tasks", &goyai.LocalizeConfig{PluralCount: 0}))

// output "There is 1 task left."
fmt.Println(i18n.Localize("en", "remaining_tasks", &goyai.LocalizeConfig{PluralCount: 1}))

// output "There are 2 tasks left."
fmt.Println(i18n.Localize("en", "remaining_tasks", &goyai.LocalizeConfig{PluralCount: 2}))

// all these commands output "Hmmm!"
fmt.Println(i18n.Localize("en", "remaining_tasks")) // no PluralCount specified, plural form "other" is used
fmt.Println(i18n.Localize("en", "remaining_tasks", &goyai.LocalizeConfig{PluralCount: 3}))  // plural form "other" is used
fmt.Println(i18n.Localize("en", "remaining_tasks", &goyai.LocalizeConfig{PluralCount: -1})) // plural form "other" is used
```

**Plural form**

A localized message can have several plural forms, specified by `zero`, `one`, `two`, `few`, `many` and `other` sections in the language file.
Plural form of a message is picked up based on the following rules:
- If `PluralCount=0` the plural form `zero` is picked up. If the message text is empty then the plural form `other` is used.
- If `PluralCount=1` the plural form `one` is picked up. If the message text is empty then the plural form `few` is chosen. If the message is, again, empty then the plural form `other` is used.
- If `PluralCount=2` the plural form `two` is picked up. If the message text is empty then the plural form `many` is chosen. If the message is, again, empty then the plural form `other` is picked up.
- Other cases (including the case when `PluralCount` is not specified), the plural form `other` is picked up.

If a message is defined by a simple string (e.g. `hello: Hello, world!`), the string is the content of the message's plural form `other` and all other plural forms are empty.

## Contributing

Use [Github issues](https://github.com/btnguyen2k/goyai/issues) for bug reports and feature requests.

Contribute by Pull Request:

1. Fork `goyai` on github (https://help.github.com/articles/fork-a-repo/)
2. Create a topic branch (`git checkout -b my_branch`)
3. Implement your change
4. Push to your branch (`git push origin my_branch`)
5. Post a pull request on github (https://help.github.com/articles/creating-a-pull-request/)

## License

MIT - see [LICENSE.md](LICENSE.md).
