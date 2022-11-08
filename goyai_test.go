package goyai

import (
	"os"
	"testing"
)

func TestNullI18n(t *testing.T) {
	testName := "TestNullI18n"
	i18n := NullI18n()
	if i18n == nil {
		t.Fatalf("%s failed: nil", testName)
	}
	if e, v := "", i18n.Localize("", ""); e != v {
		t.Fatalf("%s failed", testName)
	}
}

const (
	msgIdSimple          = "hello"
	msgTextSimple        = "Hello, world"
	msgIdSimpleWho       = "hello_who"
	msgTextSimpleWho     = "Hello Thanh"
	msgIdSimpleArbitrary = "hello_arbitrary"
)

const (
	_other = "Other cases"
	_zero  = "There is no item"
	_one   = "There is one item"
	_two   = "There is two items"
)

const jsonContent = `{
  "en": {
    "_name": "English",
    "hello": "Hello, world",
	"hello_who": "Hello {{.name}}",
	"hello_arbitrary": "Hello {{.name|upper}}",
    "count": {
        "desc": "Demo plural forms",
        "zero": "There is no item",
        "One": "There is one item",
        "TWO": "There is two items",
        "Other": "Other cases"
    }
  },
  "en2": {},
  "t1": {}
}`
const jsonFile = "test_all_in_one.json"

const yamlContent = `---
en:
en2:
  _name: "English"
  hello: "Hello, world"
  hello_who: "Hello {{.name}}"
  hello_arbitrary: "Hello {{.name}}"
  count:
    "desc": "Demo plural forms"
    "zero": "There is no item"
    "One": "There is one item"
    "TWO": "There is two items"
    "Other": "Other cases"
  other: "others"
t2:
`
const yamlFile = "test_all_in_one.yaml"

const tempDir = "temp/"

func _initDataJson() error {
	os.Mkdir(tempDir, 0711)
	f, err := os.Create(tempDir + jsonFile)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.WriteString(jsonContent); err != nil {
		return err
	}
	return nil
}

func _initDataYaml() error {
	os.Mkdir(tempDir, 0711)
	f, err := os.Create(tempDir + yamlFile)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.WriteString(yamlContent); err != nil {
		return err
	}
	return nil
}

func TestBuildI18n_InvalidFormat(t *testing.T) {
	testName := "TestBuildI18n_InvalidFormat"
	i18n, err := BuildI18n(I18nOptions{I18nFileFormat: Auto - 1})
	if i18n != nil || err == nil {
		t.Fatalf("%s failed", testName)
	}
}

func TestBuildI18n_FileNotExists(t *testing.T) {
	testName := "TestBuildI18n_FileNotExists"
	for _, format := range []I18nFileFormat{Auto, Json, Yaml} {
		i18n, err := BuildI18n(I18nOptions{I18nFileFormat: format, ConfigFileOrDir: "not-exists"})
		if i18n != nil || err == nil {
			t.Fatalf("%s failed", testName)
		}
	}
}

func TestBuildI18n_InvalidFileFormat(t *testing.T) {
	testName := "TestBuildI18n_InvalidFileFormat"

	os.RemoveAll(tempDir)
	os.Mkdir(tempDir, 0711)
	{
		f, err := os.Create(tempDir + jsonFile)
		if err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
		defer f.Close()
		if _, err = f.WriteString(yamlContent); err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
	}
	{
		f, err := os.Create(tempDir + yamlFile)
		if err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
		defer f.Close()
		if _, err = f.WriteString(jsonContent[1:]); err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
	}

	if i18n, err := BuildI18n(I18nOptions{ConfigFileOrDir: tempDir, I18nFileFormat: Json}); i18n != nil || err == nil {
		t.Fatalf("%s failed (format=Json)", testName)
	}
	if i18n, err := BuildI18n(I18nOptions{ConfigFileOrDir: tempDir, I18nFileFormat: Yaml}); i18n != nil || err == nil {
		t.Fatalf("%s failed (format=Yaml)", testName)
	}
	if i18n, err := BuildI18n(I18nOptions{ConfigFileOrDir: tempDir, I18nFileFormat: Auto}); i18n != nil || err == nil {
		t.Fatalf("%s failed (format=Auto)", testName)
	}
}

func TestBuildI18n_SingleFile_JsonAuto(t *testing.T) {
	testName := "TestBuildI18n_SingleFile_JsonAuto"

	os.RemoveAll(tempDir)
	_initDataJson()
	{
		i18n, err := BuildI18n(I18nOptions{ConfigFileOrDir: tempDir + jsonFile, I18nFileFormat: Json})
		if i18n == nil || err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
		if e, v := 3, len(i18n.AvailableLocales()); v != e {
			t.Fatalf("%s failed, expected %d available locales but received %d", testName, e, v)
		}
	}
	{
		i18n, err := BuildI18n(I18nOptions{ConfigFileOrDir: tempDir + jsonFile, I18nFileFormat: Auto})
		if i18n == nil || err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
		if e, v := 3, len(i18n.AvailableLocales()); v != e {
			t.Fatalf("%s failed, expected %d available locales but received %d", testName, e, v)
		}
	}
}

func TestBuildI18n_SingleFile_YamlAuto(t *testing.T) {
	testName := "TestBuildI18n_SingleFile_YamlAuto"

	os.RemoveAll(tempDir)
	_initDataYaml()
	{
		i18n, err := BuildI18n(I18nOptions{ConfigFileOrDir: tempDir + yamlFile, I18nFileFormat: Yaml})
		if i18n == nil || err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
		if e, v := 3, len(i18n.AvailableLocales()); v != e {
			t.Fatalf("%s failed, expected %d available locales but received %d", testName, e, v)
		}
	}
	{
		i18n, err := BuildI18n(I18nOptions{ConfigFileOrDir: tempDir + yamlFile, I18nFileFormat: Auto})
		if i18n == nil || err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
		if e, v := 3, len(i18n.AvailableLocales()); v != e {
			t.Fatalf("%s failed, expected %d available locales but received %d", testName, e, v)
		}
	}
}

func TestBuildI18n_Directory(t *testing.T) {
	testName := "TestBuildI18n_Directory"

	os.RemoveAll(tempDir)
	_initDataJson()
	_initDataYaml()
	{
		i18n, err := BuildI18n(I18nOptions{ConfigFileOrDir: tempDir, I18nFileFormat: Json})
		if i18n == nil || err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
		if e, v := 3, len(i18n.AvailableLocales()); v != e {
			t.Fatalf("%s failed, expected %d available locales but received %d", testName, e, v)
		}
	}
	{
		i18n, err := BuildI18n(I18nOptions{ConfigFileOrDir: tempDir, I18nFileFormat: Yaml})
		if i18n == nil || err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
		if e, v := 3, len(i18n.AvailableLocales()); v != e {
			t.Fatalf("%s failed, expected %d available locales but received %d", testName, e, v)
		}
	}
	{
		i18n, err := BuildI18n(I18nOptions{ConfigFileOrDir: tempDir, I18nFileFormat: Auto})
		if i18n == nil || err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
		if e, v := 4, len(i18n.AvailableLocales()); v != e {
			t.Fatalf("%s failed, expected %d available locales but received %d", testName, e, v)
		}
	}
}

func TestGoi18n_Localize_Simple(t *testing.T) {
	testName := "TestGoi18n_Localize_Simple"

	os.RemoveAll(tempDir)
	_initDataJson()
	i18n, err := BuildI18n(I18nOptions{ConfigFileOrDir: tempDir + jsonFile, I18nFileFormat: Auto})
	if i18n == nil || err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if e, v := msgTextSimple, i18n.Localize("en", msgIdSimple); v != e {
		t.Fatalf("%s failed: msg-id [%s] / expected [%s] but received %s", testName, msgIdSimple, e, v)
	}
	msgDefault := "default message"
	if e, v := msgDefault, i18n.Localise("en", msgIdSimple+"-notfound", &LocalizeConfig{DefaultMessage: msgDefault}); v != e {
		t.Fatalf("%s failed: msg-id [%s] / expected [%s] but received %s", testName, msgIdSimple, e, v)
	}
}

func TestGoi18n_Localize_Simple_DefaultLocale(t *testing.T) {
	testName := "TestGoi18n_Localize_Simple_DefaultLocale"

	os.RemoveAll(tempDir)
	_initDataYaml()
	i18n, err := BuildI18n(I18nOptions{ConfigFileOrDir: tempDir + yamlFile, I18nFileFormat: Auto, DefaultLocale: "en2"})
	if i18n == nil || err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if e, v := msgTextSimple, i18n.Localize("notfound", msgIdSimple); v != e {
		t.Fatalf("%s failed: msg-id [%s] / expected [%s] but received %s", testName, msgIdSimple, e, v)
	}
	msgDefault := "default message"
	if e, v := msgDefault, i18n.Localise("notfound", msgIdSimple+"-notfound", &LocalizeConfig{DefaultMessage: msgDefault}); v != e {
		t.Fatalf("%s failed: msg-id [%s] / expected [%s] but received %s", testName, msgIdSimple, e, v)
	}
}

func TestGoi18n_Localize_Simple_NoLocale(t *testing.T) {
	testName := "TestGoi18n_Localize_Simple_NoLocale"

	os.RemoveAll(tempDir)
	_initDataYaml()
	i18n, err := BuildI18n(I18nOptions{ConfigFileOrDir: tempDir + yamlFile, I18nFileFormat: Auto, DefaultLocale: "notfound"})
	if i18n == nil || err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if e, v := "", i18n.Localize("notexists", msgIdSimple); v != e {
		t.Fatalf("%s failed: msg-id [%s] / expected [%s] but received %s", testName, msgIdSimple, e, v)
	}
	msgDefault := "default message"
	if e, v := msgDefault, i18n.Localise("notexists", msgIdSimple+"-notfound", LocalizeConfig{DefaultMessage: msgDefault}); v != e {
		t.Fatalf("%s failed: msg-id [%s] / expected [%s] but received %s", testName, msgIdSimple, e, v)
	}
}

func TestGoi18n_LocalizeWithParam_Simple(t *testing.T) {
	testName := "TestGoi18n_LocalizeWithParam_Simple"

	os.RemoveAll(tempDir)
	_initDataJson()
	i18n, err := BuildI18n(I18nOptions{ConfigFileOrDir: tempDir + jsonFile, I18nFileFormat: Auto})
	if i18n == nil || err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if e, v := msgTextSimpleWho, i18n.Localize("en", msgIdSimpleWho, &LocalizeConfig{TemplateData: map[string]interface{}{"name": "Thanh"}}); v != e {
		t.Fatalf("%s failed: msg-id [%s] / expected [%s] but received %s", testName, msgIdSimple, e, v)
	}
}

func TestGoi18n_LocalizeWithParam_Simple_DefaultLocale(t *testing.T) {
	testName := "TestGoi18n_LocalizeWithParam_Simple_DefaultLocale"

	os.RemoveAll(tempDir)
	_initDataYaml()
	i18n, err := BuildI18n(I18nOptions{ConfigFileOrDir: tempDir + yamlFile, I18nFileFormat: Auto, DefaultLocale: "en2"})
	if i18n == nil || err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if e, v := msgTextSimpleWho, i18n.Localise("notfound", msgIdSimpleWho, &LocalizeConfig{TemplateData: map[string]interface{}{"name": "Thanh"}}); v != e {
		t.Fatalf("%s failed: msg-id [%s] / expected [%s] but received %s", testName, msgIdSimple, e, v)
	}
}

func TestGoi18n_Localize_Plural(t *testing.T) {
	testName := "TestGoi18n_Localize_Plural"

	os.RemoveAll(tempDir)
	_initDataJson()
	i18n, err := BuildI18n(I18nOptions{ConfigFileOrDir: tempDir + jsonFile, I18nFileFormat: Auto})
	if i18n == nil || err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	expected := map[interface{}]string{"none": _other, -2: _other, -1: _other, 0: _zero, 1: _one, 2: _two, 3: _other, 4: _other}
	for k, e := range expected {
		v := i18n.Localize("en", "count", &LocalizeConfig{PluralCount: k})
		if v != e {
			t.Fatalf("%s failed (%v), expect [%s] but received [%s]", testName, k, e, v)
		}
	}

	if e, v := _other, i18n.Localize("en", "count", nil); v != e {
		if v != e {
			t.Fatalf("%s failed (%v), expect [%s] but received [%s]", testName, "<nil>", e, v)
		}
	}
	if e, v := _other, i18n.Localize("en", "count", &LocalizeConfig{}); v != e {
		if v != e {
			t.Fatalf("%s failed (%v), expect [%s] but received [%s]", testName, "<empty>", e, v)
		}
	}
}

func TestGoi18n_Localize_Arbitrary(t *testing.T) {
	testName := "TestGoi18n_Localize_Arbitrary"

	os.RemoveAll(tempDir)
	_initDataYaml()
	i18n, err := BuildI18n(I18nOptions{ConfigFileOrDir: tempDir + yamlFile, I18nFileFormat: Auto})
	if i18n == nil || err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	data := map[interface{}]string{
		"A String":    "Hello A String",
		false:         "Hello false",
		true:          "Hello true",
		0:             "Hello 0",
		int(-1):       "Hello -1",
		float32(2.3):  "Hello 2.3",
		float64(-3.4): "Hello -3.4",
		uint(5):       "Hello 5",
		int8(-6):      "Hello -6",
		uint8(7):      "Hello 7",
		int16(-8):     "Hello -8",
		uint16(9):     "Hello 9",
		int32(-10):    "Hello -10",
		uint32(11):    "Hello 11",
		int64(-12):    "Hello -12",
		uint64(13):    "Hello 13",
	}
	for param, expected := range data {
		v := i18n.Localize("en2", msgIdSimpleArbitrary, param)
		if v != expected {
			t.Fatalf("%s failed: expect [%s] but received [%s]", testName, expected, v)
		}
	}
}
