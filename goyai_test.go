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
	msgIdSimple   = "hello"
	msgTextSimple = "Hello, world"
)

const jsonContent = `{
  "en": {
    "_name": "English",
    "hello": "Hello, world",
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

	{
		os.Mkdir(tempDir, 0711)
		f, err := os.Create(tempDir + jsonFile)
		if err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
		defer f.Close()
		if _, err = f.WriteString(jsonContent); err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
	}
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

	{
		os.Mkdir(tempDir, 0711)
		f, err := os.Create(tempDir + yamlFile)
		if err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
		defer f.Close()
		if _, err = f.WriteString(yamlContent); err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
	}
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
	os.Mkdir(tempDir, 0711)
	{
		f, err := os.Create(tempDir + jsonFile)
		if err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
		defer f.Close()
		if _, err = f.WriteString(jsonContent); err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
	}
	{
		f, err := os.Create(tempDir + yamlFile)
		if err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
		defer f.Close()
		if _, err = f.WriteString(yamlContent); err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
	}

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

	{
		os.Mkdir(tempDir, 0711)
		f, err := os.Create(tempDir + jsonFile)
		if err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
		defer f.Close()
		if _, err = f.WriteString(jsonContent); err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
	}

	i18n, err := BuildI18n(I18nOptions{ConfigFileOrDir: tempDir + jsonFile, I18nFileFormat: Auto})
	if i18n == nil || err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if e, v := msgTextSimple, i18n.Localize("en", msgIdSimple); v != e {
		t.Fatalf("%s failed: msg-id [%s] / expected [%s] but received %s", testName, msgIdSimple, e, v)
	}
	msgDefault := "default message"
	if e, v := msgDefault, i18n.Localize("en", msgIdSimple+"-notfound", &LocalizeConfig{DefaultMessage: msgDefault}); v != e {
		t.Fatalf("%s failed: msg-id [%s] / expected [%s] but received %s", testName, msgIdSimple, e, v)
	}
}

func TestGoi18n_Localize_Simple_DefaultLocale(t *testing.T) {
	testName := "TestGoi18n_Localize_Simple_DefaultLocale"

	{
		os.Mkdir(tempDir, 0711)
		f, err := os.Create(tempDir + yamlFile)
		if err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
		defer f.Close()
		if _, err = f.WriteString(yamlContent); err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
	}

	i18n, err := BuildI18n(I18nOptions{ConfigFileOrDir: tempDir + yamlFile, I18nFileFormat: Auto, DefaultLocale: "en2"})
	if i18n == nil || err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if e, v := msgTextSimple, i18n.Localize("notfound", msgIdSimple); v != e {
		t.Fatalf("%s failed: msg-id [%s] / expected [%s] but received %s", testName, msgIdSimple, e, v)
	}
	msgDefault := "default message"
	if e, v := msgDefault, i18n.Localize("notfound", msgIdSimple+"-notfound", &LocalizeConfig{DefaultMessage: msgDefault}); v != e {
		t.Fatalf("%s failed: msg-id [%s] / expected [%s] but received %s", testName, msgIdSimple, e, v)
	}
}

func TestGoi18n_Localize_Simple_NoLocale(t *testing.T) {
	testName := "TestGoi18n_Localize_Simple_NoLocale"

	{
		os.Mkdir(tempDir, 0711)
		f, err := os.Create(tempDir + yamlFile)
		if err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
		defer f.Close()
		if _, err = f.WriteString(yamlContent); err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
	}

	i18n, err := BuildI18n(I18nOptions{ConfigFileOrDir: tempDir + yamlFile, I18nFileFormat: Auto, DefaultLocale: "notfound"})
	if i18n == nil || err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if e, v := "", i18n.Localize("notexists", msgIdSimple); v != e {
		t.Fatalf("%s failed: msg-id [%s] / expected [%s] but received %s", testName, msgIdSimple, e, v)
	}
	msgDefault := "default message"
	if e, v := msgDefault, i18n.Localize("notexists", msgIdSimple+"-notfound", &LocalizeConfig{DefaultMessage: msgDefault}); v != e {
		t.Fatalf("%s failed: msg-id [%s] / expected [%s] but received %s", testName, msgIdSimple, e, v)
	}
}
