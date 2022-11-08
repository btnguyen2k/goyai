package goyai

import (
	"log"
	"regexp"
	"sort"
	"sync"
)

// Goi18n is the default I18n implementation from goyai.
type Goi18n struct {
	defaultLocale string
	locales       map[string]*LocaleInfo
	cachedLocales []LocaleInfo
	messagesStore map[string]map[string]*Message // {locale->{msg-id->msg-data}}
	lock          sync.Mutex
}

// Localise implements I18n.Localise
func (i *Goi18n) Localise(locale, msgId string, params ...interface{}) string {
	return i.Localize(locale, msgId, params...)
}

func _extractFirstConfig(params ...interface{}) *LocalizeConfig {
	for _, param := range params {
		if cfg, ok := param.(LocalizeConfig); ok {
			return &cfg
		}
		if cfg, ok := param.(*LocalizeConfig); ok && cfg != nil {
			return cfg
		}
	}
	return nil
}

var rePlaceholderToken = regexp.MustCompile(`{{\$?\.([\w]+).*?}}`)

func _buildTemplateData(msg string, params ...interface{}) map[string]interface{} {
	templateData := make(map[string]interface{})
	forwardMap := make(map[string]int)
	reverseMap := make(map[int]string)
	matches := rePlaceholderToken.FindAllStringSubmatch(msg, -1)
	index := 0
	for _, match := range matches {
		token := match[1]
		if _, ok := forwardMap[token]; !ok {
			forwardMap[token] = index
			reverseMap[index] = token
			index++
		}
	}
	index = 0
	for _, param := range params {
		switch v := param.(type) {
		case string, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
			if token, ok := reverseMap[index]; ok {
				templateData[token] = v
			}
		}
		index++
	}
	return templateData
}

// Localize implements I18n.Localize
func (i *Goi18n) Localize(locale, msgId string, params ...interface{}) string {
	cfg := _extractFirstConfig(params...)
	var msg string
	localizedMessage := i.getLocalizedMessage(msgId, locale, i.defaultLocale)
	if localizedMessage != nil {
		if cfg == nil && len(params) > 0 {
			cfg = &LocalizeConfig{TemplateData: _buildTemplateData(localizedMessage.Other, params...)}
		}
		msg = localizedMessage.render(cfg)
	}
	if msg == "" {
		log.Printf("[WARN] localized message [%s] not defined for locale [%s]", msgId, locale)
	}
	if msg == "" && cfg != nil {
		msg = cfg.DefaultMessage
	}
	return msg
}

func (i *Goi18n) getLocalizedMessage(msgId, locale, defaultLocale string) *Message {
	if i.messagesStore == nil {
		return nil
	}
	if locale == "" || i.locales[locale] == nil {
		if locale != "" {
			log.Printf("[WARN] locale [%s] not exist, revert back to default", locale)
		}
		locale = defaultLocale
	}
	localizedMessagesData := i.messagesStore[locale]
	if localizedMessagesData != nil {
		return localizedMessagesData[msgId]
	}
	return nil
}

// AvailableLocales implements I18n.AvailableLocales.
func (i *Goi18n) AvailableLocales() []LocaleInfo {
	i.lock.Lock()
	defer i.lock.Unlock()

	if i.cachedLocales == nil {
		i.cachedLocales = make([]LocaleInfo, len(i.locales))
		var j = 0
		for _, localeInfo := range i.locales {
			i.cachedLocales[j] = *localeInfo
			j++
		}
		sort.Slice(i.cachedLocales, func(x, y int) bool {
			return i.cachedLocales[x].DisplayName < i.cachedLocales[y].DisplayName
		})
	}

	return i.cachedLocales
}
