// Package goyai provides a simple interface to support i18n for Go applications.
package goyai

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/btnguyen2k/consu/reddo"
	"gopkg.in/yaml.v2"
)

const (
	// Version of goyai package
	Version = "0.1.0"
)

// LocaleInfo captures info of a locale package.
type LocaleInfo struct {
	Id          string
	DisplayName string
}

// LocalizeConfig configures how a message should be localised, used by function I18n.Localize.
type LocalizeConfig struct {
	// TemplateData is used to transform the message's template.
	TemplateData map[string]interface{}

	// PluralCount determines which plural form of the message is used. PluralCount must be an integer or nil.
	PluralCount interface{}

	// DefaultMessage holds the default message where there is no localized one.
	DefaultMessage string
}

// I18n is the main interface of goyai package, offering APIs to render a message.
type I18n interface {
	// Localize returns a localized message. Config is optional, and only the first supplied config is accounted for.
	Localize(locale, msgId string, config ...*LocalizeConfig) string

	// AvailableLocales returns all defined locale configurations.
	AvailableLocales() []LocaleInfo
}

// I18nFileFormat defines list of supported i18n configuration file formats.
type I18nFileFormat int

const (
	// Auto hints that the file format should be automatically detected.
	Auto I18nFileFormat = iota

	// Json hints that the file is JSON-encoded.
	Json

	// Yaml hints that the file is YAML-encoded.
	Yaml
)

var (
	// ErrInvalidFileFormat indicates that the specified language file format is not supported.
	ErrInvalidFileFormat = errors.New("language file format is invalid or not supported")
)

// I18nOptions specifies options to build new I18n instances.
type I18nOptions struct {
	// ConfigFileOrDir points to the configuration file or the directory where configuration files are located.
	ConfigFileOrDir string

	// DefaultLocale is the default locale to be used when non specified.
	DefaultLocale string

	// I18nFileFormat hints the format of configuration files.
	I18nFileFormat I18nFileFormat
}

// NullI18n returns a "null" I18n instance.
func NullI18n() I18n {
	return &Goi18n{}
}

// BuildI18n builds an I18n instance from message file(s) and returns it.
func BuildI18n(opts I18nOptions) (I18n, error) {
	switch opts.I18nFileFormat {
	case Auto:
		return buildI18nAuto(opts)
	case Json:
		return buildI18nJson(opts)
	case Yaml:
		return buildI18nYaml(opts)
	default:
		return nil, ErrInvalidFileFormat
	}
}

func buildI18nAuto(opts I18nOptions) (I18n, error) {
	localesStore := make(map[string]*LocaleInfo)
	messagesStore := make(map[string]map[string]*Message)

	fileOrDir := opts.ConfigFileOrDir
	fileInfo, err := os.Stat(fileOrDir)
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		// a directory contains multiple language files
		files, err := ioutil.ReadDir(fileOrDir)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			normFilename := strings.ToLower(file.Name())
			if !file.IsDir() && (strings.HasSuffix(normFilename, ".yaml") || strings.HasSuffix(normFilename, ".yml")) {
				if err := loadLangFileYaml(localesStore, messagesStore, fileOrDir, file); err != nil {
					return nil, err
				}
			}
			if !file.IsDir() && strings.HasSuffix(normFilename, ".json") {
				if err := loadLangFileJson(localesStore, messagesStore, fileOrDir, file); err != nil {
					return nil, err
				}
			}
		}
	} else {
		// a single language file
		normFilename := strings.ToLower(fileInfo.Name())
		if strings.HasSuffix(normFilename, ".yaml") || strings.HasSuffix(normFilename, ".yml") {
			if err := loadLangFileYaml(localesStore, messagesStore, filepath.Dir(fileOrDir), fileInfo); err != nil {
				return nil, err
			}
		}
		if strings.HasSuffix(normFilename, ".json") {
			if err := loadLangFileJson(localesStore, messagesStore, filepath.Dir(fileOrDir), fileInfo); err != nil {
				return nil, err
			}
		}
	}

	return &Goi18n{defaultLocale: opts.DefaultLocale, locales: localesStore, messagesStore: messagesStore}, nil
}

func buildI18nJson(opts I18nOptions) (I18n, error) {
	localesStore := make(map[string]*LocaleInfo)
	messagesStore := make(map[string]map[string]*Message)

	fileOrDir := opts.ConfigFileOrDir
	fileInfo, err := os.Stat(fileOrDir)
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		// a directory contains multiple language files
		files, err := ioutil.ReadDir(fileOrDir)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".json") {
				if err := loadLangFileJson(localesStore, messagesStore, fileOrDir, file); err != nil {
					return nil, err
				}
			}
		}
	} else if err := loadLangFileJson(localesStore, messagesStore, filepath.Dir(fileOrDir), fileInfo); err != nil { // a single language file
		return nil, err
	}

	return &Goi18n{defaultLocale: opts.DefaultLocale, locales: localesStore, messagesStore: messagesStore}, nil
}

func loadLangFileJson(localesStore map[string]*LocaleInfo, messagesStore map[string]map[string]*Message, dirPath string, file os.FileInfo) error {
	buf, err := ioutil.ReadFile(dirPath + "/" + file.Name())
	if err != nil {
		return err
	}
	var langData map[string]map[string]interface{}
	if err := json.Unmarshal(buf, &langData); err != nil {
		return err
	}
	return parseLangData(localesStore, messagesStore, langData)
}

func buildI18nYaml(opts I18nOptions) (I18n, error) {
	localesStore := make(map[string]*LocaleInfo)
	messagesStore := make(map[string]map[string]*Message)

	fileOrDir := opts.ConfigFileOrDir
	fileInfo, err := os.Stat(fileOrDir)
	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		// a directory contains multiple language files
		files, err := ioutil.ReadDir(fileOrDir)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			normFilename := strings.ToLower(file.Name())
			if !file.IsDir() && (strings.HasSuffix(normFilename, ".yaml") || strings.HasSuffix(normFilename, ".yml")) {
				if err := loadLangFileYaml(localesStore, messagesStore, fileOrDir, file); err != nil {
					return nil, err
				}
			}
		}
	} else if err := loadLangFileYaml(localesStore, messagesStore, filepath.Dir(fileOrDir), fileInfo); err != nil { // a single language file
		return nil, err
	}

	return &Goi18n{defaultLocale: opts.DefaultLocale, locales: localesStore, messagesStore: messagesStore}, nil
}

func loadLangFileYaml(localesStore map[string]*LocaleInfo, messagesStore map[string]map[string]*Message, dirPath string, file os.FileInfo) error {
	buf, err := ioutil.ReadFile(dirPath + "/" + file.Name())
	if err != nil {
		return err
	}
	var langData map[string]map[string]interface{}
	if err := yaml.Unmarshal(buf, &langData); err != nil {
		return err
	}
	return parseLangData(localesStore, messagesStore, langData)
}

func parseLangData(localesStore map[string]*LocaleInfo, messagesStore map[string]map[string]*Message, langData map[string]map[string]interface{}) error {
	// top level is "locale" mapped to messages
	for locale, msgMap := range langData {
		localeInfo := localesStore[locale]
		if localeInfo == nil {
			localeInfo = &LocaleInfo{Id: locale, DisplayName: locale}
			localesStore[locale] = localeInfo
		}

		localizedMessages := messagesStore[locale]
		if localizedMessages == nil {
			localizedMessages = make(map[string]*Message)
			messagesStore[locale] = localizedMessages
		}

		for msgId, msgData := range msgMap {
			// special message-id
			if (msgId == "_display" || msgId == "_name") && (localeInfo.DisplayName == "" || localeInfo.DisplayName == localeInfo.Id) {
				localeInfo.DisplayName, _ = reddo.ToString(msgData)
				continue
			}

			// message-id mapped to message data, which is either simply a string or a struct
			if msg, err := ParseMessage(msgId, msgData); err != nil {
				return err
			} else {
				localizedMessages[msgId] = msg
			}
		}
	}

	return nil
}
