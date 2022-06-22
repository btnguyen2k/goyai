package goyai

import (
	"log"
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

// Localize implements I18n.Localize
func (i *Goi18n) Localize(locale, msgId string, config ...*LocalizeConfig) string {
	var cfg *LocalizeConfig
	if len(config) > 0 && config[0] != nil {
		cfg = config[0]
	}

	var msg string
	localizedMessage := i.getLocalizedMessage(msgId, locale, i.defaultLocale)
	if localizedMessage != nil {
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
