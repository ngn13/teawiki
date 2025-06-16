package config

import (
	"fmt"
	"net/url"
	"time"

	"github.com/ngn13/ortam"
)

type Config struct {
	Debug      bool
	ListenAddr string `ortam:"LISTEN"`

	Url      *url.URL
	Name     string
	Desc     string
	Keywords string

	SourceUrl *url.URL
	CommitUrl *url.URL

	RepoUrl  *url.URL
	RepoPath string

	ReloadInterval time.Duration
	WebhookSecret  string

	ChromaStyle string `ortam:"CHROMA"`
	Theme       string
	Lang        string
	Time        string

	Logo string
	Icon string
}

func (c *Config) TimeStr(t time.Time) string {
	return t.Format(c.Time)
}

func Load() (*Config, error) {
	config := Config{
		// default options
		Debug:      false,
		ListenAddr: "127.0.0.1:8080",

		Url:      nil,
		Name:     "my wiki",
		Desc:     "my personal wiki",
		Keywords: "wiki",

		SourceUrl: nil,
		CommitUrl: nil,

		RepoPath: "./source",
		RepoUrl:  nil,

		ReloadInterval: time.Minute * 30,
		WebhookSecret:  "",

		ChromaStyle: "",
		Theme:       "dark",
		Lang:        "en",
		Time:        "02/01/06 15:04:05 GMT-07",

		Logo: "logo.png",
		Icon: "logo.png",
	}

	if err := ortam.Load(&config, "TW"); err != nil {
		return nil, err
	}

	if config.Name == "" || config.Desc == "" || config.ListenAddr == "" ||
		config.Theme == "" || config.RepoPath == "" {
		return nil, fmt.Errorf(
			"a required configuration option is missing, please see docs/config.md",
		)
	}

	if config.ReloadInterval == 0 {
		return nil, fmt.Errorf("reload interval should be non-zero")
	}

	if config.ChromaStyle == "" {
		switch config.Theme {
		case "dark":
			config.ChromaStyle = "rrt"

		case "light":
			config.ChromaStyle = "emacs"
		}
	}

	return &config, nil
}
