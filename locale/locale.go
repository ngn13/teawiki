package locale

import (
	"fmt"
	"strconv"

	"github.com/goccy/go-yaml"
	"github.com/ngn13/teawiki/config"
	"github.com/ngn13/teawiki/log"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const (
	LOCALE_DIR     = "locale"
	LOCALE_EXT     = "yaml"
	LOCALE_DEFAULT = LOCALE_DIR + "/en." + LOCALE_EXT
)

type Locale struct {
	Conf      *config.Config
	Bundle    *i18n.Bundle
	Localizer *i18n.Localizer
}

func (l *Locale) Get(id string, args ...interface{}) string {
	var (
		data   map[string]interface{} = nil
		plural int                    = 2
	)

	for i, arg := range args {
		if data == nil {
			data = make(map[string]interface{})
		}

		data["a"+strconv.Itoa(i)] = arg

		if iarg, ok := arg.(int); i == 0 && ok && iarg == 1 {
			plural = 1
		}
	}

	str, err := l.Localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: data,
		PluralCount:  plural,
	})

	if err != nil {
		log.Warn("failed to get the locale for %s: %s", id, err.Error())
		return id
	}

	return str
}

func New(conf *config.Config) (*Locale, error) {
	var (
		locale_path string = fmt.Sprintf(
			"%s/%s.%s",
			LOCALE_DIR,
			conf.Lang,
			LOCALE_EXT,
		)
		err error
	)

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc(LOCALE_EXT, yaml.Unmarshal)

	if _, err = bundle.LoadMessageFile(LOCALE_DEFAULT); err != nil {
		return nil, err
	}

	if _, err = bundle.LoadMessageFile(locale_path); err != nil {
		return nil, err
	}

	return &Locale{
		Conf:      conf,
		Bundle:    bundle,
		Localizer: i18n.NewLocalizer(bundle, conf.Lang),
	}, nil
}
