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
- Can be used in/integrated with [html/template](http://golang.org/pkg/html/template/) (since [v0.2.0](RELEASE-NOTES.md)).

## Installation

```go
go get github.com/btnguyen2k/goyai
```

## Usage & Documentation

[![PkgGoDev](https://pkg.go.dev/badge/github.com/btnguyen2k/goyai)](https://pkg.go.dev/github.com/btnguyen2k/goyai)

**Language file format**

`goyai` supports language files in JSON and YAML (1.2) formats. Here is an example of language file in YAML:

```yaml
---
en:
  _name: English
  hello: Hello, world!
  hello_param: Hello buddy {{.name}}
  remaining_tasks:
    zero: Congratulation {{.who}}! You are free now.
    one: There is 1 task left.
    two: There are 2 tasks left.
    few: There are a few tasks left.
    many: There are many tasks left.
    other: Hmmm!

vi:
  _name: Tiếng Việt
  hello: Xin chào
  hello_param: Chào bạn {{.name}}
  remaining_tasks:
    zero: Chúc mừng {{.who}}! Bạn đã hoàn thành công việc.
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
i18n, err := BuildI18n(goyai.I18nOptions{ConfigFileOrDir: "./languages/", I18nFileFormat: goyai.Auto, DefaultLocale: "en"})
```

**Localize messages via I18n instance**

```go
// output "Xin chào"
fmt.Println(i18n.Localize("vi", "hello"))

// locale "nf" does not exist, fall back to default locale "en"
// output "Hello, world!"
fmt.Println(i18n.Localize("nf", "hello"))

// Pass template data, output "Hello buddy Thanh"
fmt.Println(i18n.Localize("en", "hello_param", goyai.LocalizeConfig{TemplateData: map[string]interface{}{"name": "Thanh"}}))

// Plural form, output "There is 1 task left."
fmt.Println(i18n.Localize("en", "remaining_tasks", goyai.LocalizeConfig{PluralCount: 1}))

// Plural form, output "There are 2 tasks left."
fmt.Println(i18n.Localize("en", "remaining_tasks", goyai.LocalizeConfig{PluralCount: 2}))

// Plural form, these commands output "Hmmm!"
fmt.Println(i18n.Localize("en", "remaining_tasks")) // no PluralCount specified, plural form "other" is used
fmt.Println(i18n.Localize("en", "remaining_tasks", goyai.LocalizeConfig{PluralCount: -1})) // plural form "other" is used

// Plural form, output "There are many tasks left."
fmt.Println(i18n.Localize("en", "remaining_tasks", goyai.LocalizeConfig{PluralCount: 3}))

// Plural form & pass template data, output "Congratulation btnguyen2k! You are free now."
fmt.Println(i18n.Localize("en", "remaining_tasks", goyai.LocalizeConfig{PluralCount: 0, TemplateData: map[string]interface{}{"name": "btnguyen2k"}}))
```

**Plural forms**

A localized message can have several plural forms, specified by `zero`, `one`, `two`, `few`, `many` and `other` attributes in the language file.
Plural form of a message is picked up based on the following rules:
- if `PluralCount` is negative number, `nil` or not cast-able to integer, the `other` form is chosen.
- if `PluralCount=0`, the `zero` form is chosen.
- if `PluralCount=1`, one of `one`/`few`/`other` forms is chosen, priority is from left to right (e.g. `one` form has the highest priority, if absent, the next one is checked)
- if `PluralCount=2`, one of `two`/`many`/`other` forms is chosen, priority is from left to right (e.g. `two` form has the highest priority, if absent, the next one is checked)
- if `PluralCount>2`, one of `many`/`other` forms is chosen, priority is from left to right (e.g. `many` form has the highest priority, if absent, the next one is checked)

If a message is defined by a simple string (e.g. `hello: Hello, world!`), the string is the content of the message's plural form `other` and all other plural forms are empty.

**Used in `html/template` template**

> `html/template` support requires [v0.2.0](RELEASE-NOTES.md) or higher.

Assuming the `I18n` instance is pass to a `html/template` template as a model named `i18n`. Then
- Template `The message: {{.i18n.Localize "en" "hello"}}` will be rendered as `The message: Hello, world!`.
- Template `The message: {{.i18n.Localize "en" "hello_param" "Thanh"}}` will be rendered as `The message: Hello buddy Thanh`.

> Plural form is current not supported if used in `html/template` template.

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
