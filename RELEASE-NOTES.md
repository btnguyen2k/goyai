# goyai release notes

## 2022-11-08 - v0.2.0

- Add function `I18n.Localise` which is alias of `I18n.Localize`.
- (Possible breaking change) Change signature of function `I18n.Localize`:
  - Old `Localize(locale, msgId string, config ...*LocalizeConfig) string`.
  - New `Localize(locale, msgId string, params ...interface{}) string`.
- Migrate yaml package to `gopkg.in/yaml.v3`.

## 2022-06-22 - v0.1.0

- Support localized text messages, with plural forms.
- Support template string with named variables following [text/template](http://golang.org/pkg/text/template/) syntax.
- Support message files in JSON and YAML formats.
