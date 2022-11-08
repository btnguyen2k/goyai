package goyai

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"strings"
	"text/template"

	"github.com/btnguyen2k/consu/reddo"
)

// ParseMessage parses data and returns a new Message instance.
//
// data must be either a string, or a map[string]string. If data is a string, the Message is constructed
// with data is the value of the plural form "other".
func ParseMessage(id string, data interface{}) (*Message, error) {
	msg := &Message{Id: id}
	if err := msg.parse(data); err != nil {
		return nil, err
	}
	return msg, nil
}

// Message represents a message that can be localized.
type Message struct {
	// Id is the message's unique identity.
	Id string

	// Description provides additional information about the message.
	Description string

	// Zero is the message's content for the CLDR plural form "zero".
	Zero string

	// One is the message's content for the CLDR plural form "one".
	One string

	// Two is the message's content for the CLDR plural form "two".
	Two string

	// Few is the message's content for the CLDR plural form "few".
	Few string

	// Many is the message's content for the CLDR plural form "many".
	Many string

	// Other is the message's content for the CLDR plural form "other".
	Other string
}

func (m *Message) parseMessageAttr(k string, v interface{}) error {
	var ok bool
	temp := strings.TrimSpace(strings.ToLower(k))
	switch temp {
	case "desc", "description":
		if m.Description, ok = v.(string); !ok {
			return fmt.Errorf("error parsing message data at '%s.%s'", m.Id, k)
		}
	case "zero":
		if m.Zero, ok = v.(string); !ok {
			return fmt.Errorf("error parsing message data at '%s.%s'", m.Id, k)
		}
	case "one":
		if m.One, ok = v.(string); !ok {
			return fmt.Errorf("error parsing message data at '%s.%s'", m.Id, k)
		}
	case "two":
		if m.Two, ok = v.(string); !ok {
			return fmt.Errorf("error parsing message data at '%s.%s'", m.Id, k)
		}
	case "few":
		if m.Few, ok = v.(string); !ok {
			return fmt.Errorf("error parsing message data at '%s.%s'", m.Id, k)
		}
	case "many":
		if m.Many, ok = v.(string); !ok {
			return fmt.Errorf("error parsing message data at '%s.%s'", m.Id, k)
		}
	case "other":
		if m.Other, ok = v.(string); !ok {
			return fmt.Errorf("error parsing message data at '%s.%s'", m.Id, k)
		}
	default:
		return fmt.Errorf("error parsing message data at '%s.%s'", m.Id, k)
	}
	return nil
}

// parse builds message info from data.
//
// See function ParseMessage for detailed format of data.
func (m *Message) parse(data interface{}) error {
	if data == nil {
		return fmt.Errorf("error parsing message data '%s' (type %T)", m.Id, data)
	}
	switch reflect.TypeOf(data).Kind() {
	case reflect.String:
		// message data is a simple string: it is the localized message itself
		m.Other = strings.TrimSpace(data.(string))
		return nil
	case reflect.Map:
		it := reflect.ValueOf(data).MapRange()
		for it.Next() {
			k, _ := it.Key().Interface().(string)
			v := it.Value().Interface()
			if err := m.parseMessageAttr(k, v); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("error parsing message data '%s' (type %T)", m.Id, data)
	}
}

func (m *Message) pluralFormTemplate(cfg *LocalizeConfig) string {
	var err error
	var pluralForm int64 = -1
	if cfg != nil && cfg.PluralCount != nil {
		if pluralForm, err = reddo.ToInt(cfg.PluralCount); err != nil {
			pluralForm = -1
		}
	}
	switch pluralForm {
	case 0:
		if m.Zero != "" {
			return m.Zero
		}
		return m.Other
	case 1:
		if m.One != "" {
			return m.One
		}
		if m.Few != "" {
			return m.Few
		}
		return m.Other
	case 2:
		if m.Two != "" {
			return m.Two
		}
		if m.Many != "" {
			return m.Many
		}
		return m.Other
	default:
		if pluralForm > 0 {
			if m.Many != "" {
				return m.Many
			}
		}
		return m.Other
	}
}

func (m *Message) render(cfg *LocalizeConfig) string {
	msg := m.pluralFormTemplate(cfg)
	t := template.New(m.Id)
	if _, err := t.Parse(msg); err != nil {
		log.Printf("[WARN] error parsing message [%s]: %s", m.Id, err)
		return msg
	}
	w := bytes.NewBufferString("")
	var templateData interface{}
	if cfg != nil {
		templateData = cfg.TemplateData
	}
	if err := t.Execute(w, templateData); err != nil {
		log.Printf("[WARN] error rendering message [%s]: %s", m.Id, err)
		return msg
	}
	return w.String()
}
