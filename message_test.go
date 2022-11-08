package goyai

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseMessageNil(t *testing.T) {
	testName := "TestParseMessageNil"
	msg, err := ParseMessage("msgid", nil)
	if msg != nil || err == nil {
		t.Fatalf("%s failed", testName)
	}
}

func TestParseMessageString(t *testing.T) {
	testName := "TestParseMessageString"
	msgId := "mid"
	strOther := "Localized message, plural form 'other'."
	msg, err := ParseMessage(msgId, strOther)
	if msg == nil || err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}

	emsg := Message{
		Id:          msgId,
		Description: "",
		Zero:        "",
		One:         "",
		Two:         "",
		Few:         "",
		Many:        "",
		Other:       strOther,
	}
	if !reflect.DeepEqual(emsg, *msg) {
		t.Fatalf("%s failed, expect\n%#v\nbut received\n%#v", testName, emsg, *msg)
	}
}

func TestParseMessageMapStringString(t *testing.T) {
	testName := "TestParseMessageMapStringString"
	msgId := "mid"
	data := map[string]string{
		"desc":   "Message description",
		"Zero":   "Localized message, plural form 'zero'",
		"one":    "Localized message, plural form 'one'",
		"TWO":    "Localized message, plural form 'two'",
		"feW":    "Localized message, plural form 'few'",
		" MaNy":  "Localized message, plural form 'many'",
		"oTHEr ": "Localized message, plural form 'other'",
	}
	msg, err := ParseMessage(msgId, data)
	if msg == nil || err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}

	emsg := Message{
		Id:          msgId,
		Description: "Message description",
		Zero:        "Localized message, plural form 'zero'",
		One:         "Localized message, plural form 'one'",
		Two:         "Localized message, plural form 'two'",
		Few:         "Localized message, plural form 'few'",
		Many:        "Localized message, plural form 'many'",
		Other:       "Localized message, plural form 'other'",
	}
	if !reflect.DeepEqual(emsg, *msg) {
		t.Fatalf("%s failed, expect\n%#v\nbut received\n%#v", testName, emsg, *msg)
	}
}

func TestParseMessageMapStringInterface(t *testing.T) {
	testName := "TestParseMessageMapStringInterface"
	msgId := "mid"
	data := map[string]interface{}{
		" Description ": "Message description",
		"Zero":          "Localized message, plural form 'zero'",
		"one":           "Localized message, plural form 'one'",
		"TWO":           "Localized message, plural form 'two'",
		"feW":           "Localized message, plural form 'few'",
		" MaNy":         "Localized message, plural form 'many'",
		"oTHEr ":        "Localized message, plural form 'other'",
	}
	msg, err := ParseMessage(msgId, data)
	if msg == nil || err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}

	emsg := Message{
		Id:          msgId,
		Description: "Message description",
		Zero:        "Localized message, plural form 'zero'",
		One:         "Localized message, plural form 'one'",
		Two:         "Localized message, plural form 'two'",
		Few:         "Localized message, plural form 'few'",
		Many:        "Localized message, plural form 'many'",
		Other:       "Localized message, plural form 'other'",
	}
	if !reflect.DeepEqual(emsg, *msg) {
		t.Fatalf("%s failed, expect\n%#v\nbut received\n%#v", testName, emsg, *msg)
	}
}

func TestParseMessageError(t *testing.T) {
	testName := "TestParseMessageError"
	msg, err := ParseMessage("msgid", 0)
	if msg != nil || err == nil {
		t.Fatalf("%s failed", testName)
	}
}

func TestParseMessageError2(t *testing.T) {
	testName := "TestParseMessageMapStringInterface"
	fields := []string{"Description", "Zero", "One", "Two", "Few", "Many", "Other", "Invalid"}
	for _, f := range fields {
		data := map[string]interface{}{f: 0}
		msg, err := ParseMessage("msgid", data)
		if msg != nil || err == nil {
			t.Fatalf("%s failed (field [%s])", testName, f)
		}
	}
}

const (
	zero  = "Localized message, plural form 'zero': {{.data}}"
	one   = "Localized message, plural form 'one': {{.data}}"
	two   = "Localized message, plural form 'two': {{.data}}"
	few   = "Localized message, plural form 'few': {{.data}}"
	many  = "Localized message, plural form 'many': {{.data}}"
	other = "Localized message, plural form 'other': {{.data}}"
)

func TestMessage_pluralForm(t *testing.T) {
	testName := "TestMessage_pluralForm"
	msgId := "mid"
	data := map[string]interface{}{"Zero": zero, "One": one, "Two": two, "Few": few, "Many": many, "Other": other}
	msg, err := ParseMessage(msgId, data)
	if msg == nil || err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	expected := map[interface{}]string{"none": other, -2: other, -1: other, 0: zero, 1: one, 2: two, 3: many, 4: many}
	for k, e := range expected {
		v := msg.pluralFormTemplate(&LocalizeConfig{PluralCount: k})
		if v != e {
			t.Fatalf("%s failed (%v), expect [%s] but received [%s]", testName, k, e, v)
		}
	}

	if e, v := other, msg.pluralFormTemplate(nil); v != e {
		if v != e {
			t.Fatalf("%s failed (%v), expect [%s] but received [%s]", testName, "<nil>", e, v)
		}
	}
	if e, v := other, msg.pluralFormTemplate(&LocalizeConfig{}); v != e {
		if v != e {
			t.Fatalf("%s failed (%v), expect [%s] but received [%s]", testName, "<empty>", e, v)
		}
	}
}

func TestMessage_pluralForm_ZeroAsOther(t *testing.T) {
	testName := "TestMessage_pluralForm_ZeroAsOther"
	msgId := "mid"
	data := map[string]interface{}{"Other": other}
	msg, err := ParseMessage(msgId, data)
	if msg == nil || err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if e, v := other, msg.pluralFormTemplate(&LocalizeConfig{PluralCount: 0}); v != e {
		t.Fatalf("%s failed (%v), expect [%s] but received [%s]", testName, 0, e, v)
	}
}

func TestMessage_pluralForm_OneAsFewOrOther(t *testing.T) {
	testName := "TestMessage_pluralForm_OneAsFewOrOther"
	msgId := "mid"
	{
		// "few" takes priority over "other" when "one" is absent
		data := map[string]interface{}{"Few": few, "Other": other}
		msg, err := ParseMessage(msgId, data)
		if msg == nil || err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
		if e, v := few, msg.pluralFormTemplate(&LocalizeConfig{PluralCount: 1}); v != e {
			t.Fatalf("%s failed (%v), expect [%s] but received [%s]", testName, 1, e, v)
		}
	}
	{
		data := map[string]interface{}{"Other": other}
		msg, err := ParseMessage(msgId, data)
		if msg == nil || err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
		if e, v := other, msg.pluralFormTemplate(&LocalizeConfig{PluralCount: 1}); v != e {
			t.Fatalf("%s failed (%v), expect [%s] but received [%s]", testName, 1, e, v)
		}
	}
}

func TestMessage_pluralForm_TwoAsManyOrOther(t *testing.T) {
	testName := "TestMessage_pluralForm_TwoAsManyOrOther"
	msgId := "mid"
	{
		// "many" takes priority over "other" when "two" is absent
		data := map[string]interface{}{"Many": many, "Other": other}
		msg, err := ParseMessage(msgId, data)
		if msg == nil || err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
		if e, v := many, msg.pluralFormTemplate(&LocalizeConfig{PluralCount: 2}); v != e {
			t.Fatalf("%s failed (%v), expect [%s] but received [%s]", testName, 2, e, v)
		}
	}
	{
		data := map[string]interface{}{"Other": other}
		msg, err := ParseMessage(msgId, data)
		if msg == nil || err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
		if e, v := other, msg.pluralFormTemplate(&LocalizeConfig{PluralCount: 2}); v != e {
			t.Fatalf("%s failed (%v), expect [%s] but received [%s]", testName, 2, e, v)
		}
	}
}

func TestMessage_pluralForm_GreaterThanTwoAsManyOrOther(t *testing.T) {
	testName := "TestMessage_pluralForm_GreaterThanTwoAsManyOrOther"
	msgId := "mid"
	{
		// "many" takes priority over "other"
		data := map[string]interface{}{"Many": many, "Other": other}
		msg, err := ParseMessage(msgId, data)
		if msg == nil || err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
		if e, v := many, msg.pluralFormTemplate(&LocalizeConfig{PluralCount: 3}); v != e {
			t.Fatalf("%s failed (%v), expect [%s] but received [%s]", testName, 2, e, v)
		}
	}
	{
		data := map[string]interface{}{"Other": other}
		msg, err := ParseMessage(msgId, data)
		if msg == nil || err != nil {
			t.Fatalf("%s failed: %s", testName, err)
		}
		if e, v := other, msg.pluralFormTemplate(&LocalizeConfig{PluralCount: 3}); v != e {
			t.Fatalf("%s failed (%v), expect [%s] but received [%s]", testName, 2, e, v)
		}
	}
}

func TestMessage_render(t *testing.T) {
	testName := "TestMessage_render"
	msgId := "mid"
	data := map[string]interface{}{"Zero": zero, "One": one, "Two": two, "Few": few, "Many": many, "Other": other}
	msg, err := ParseMessage(msgId, data)
	if msg == nil || err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	cfg := &LocalizeConfig{TemplateData: map[string]interface{}{"data": "value"}}
	expected := map[int]string{-2: other, -1: other, 0: zero, 1: one, 2: two, 3: many, 4: many}
	for k, _e := range expected {
		cfg.PluralCount = k
		v := msg.render(cfg)
		e := strings.ReplaceAll(_e, "{{.data}}", "value")
		if v != e {
			t.Fatalf("%s failed (%v), expect [%s] but received [%s]", testName, k, e, v)
		}
	}
}

func TestMessage_render_empty(t *testing.T) {
	testName := "TestMessage_render_empty"
	msgId := "mid"
	data := map[string]interface{}{}
	msg, err := ParseMessage(msgId, data)
	if msg == nil || err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if e, v := "", msg.render(nil); v != e {
		t.Fatalf("%s failed, expect [%s] but received [%s]", testName, e, v)
	}
}

func TestMessage_render_error(t *testing.T) {
	testName := "TestMessage_render_error"
	msgId := "mid"
	invalidTemplate := "This template is {{.invalid>"
	validTemplate := "This template is {{if .valid gt 0}}valid{{end}}"
	data := map[string]interface{}{"other": invalidTemplate, "zero": validTemplate}
	msg, err := ParseMessage(msgId, data)
	if msg == nil || err != nil {
		t.Fatalf("%s failed: %s", testName, err)
	}
	if e, v := invalidTemplate, msg.render(nil); v != e {
		t.Fatalf("%s failed, expect [%s] but received [%s]", testName, e, v)
	}
	if e, v := validTemplate, msg.render(&LocalizeConfig{PluralCount: 0}); v != e {
		t.Fatalf("%s failed, expect [%s] but received [%s]", testName, e, v)
	}
}
