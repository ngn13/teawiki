package main

import (
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/ngn13/teawiki/config"
	"github.com/ngn13/teawiki/locale"
	"github.com/ngn13/teawiki/log"
	"github.com/ngn13/teawiki/repo"
	"github.com/ngn13/teawiki/routes"
	"github.com/ngn13/teawiki/util"
)

func main() {
	var (
		conf *config.Config
		loc  *locale.Locale
		rep  *repo.Repo
		app  *fiber.App
		err  error
	)

	log.Banner()

	if conf, err = config.Load(); err != nil {
		log.Fail("failed to load the configuration: %s", err.Error())
		os.Exit(1)
	}

	if loc, err = locale.New(conf); err != nil {
		log.Fail("failed to load the locale: %s", err.Error())
		os.Exit(1)
	}

	if rep, err = repo.New(conf, loc); err != nil {
		log.Fail("failed to load the git repo: %s", err.Error())
		os.Exit(1)
	}

	// setup channels for background thread
	signal_chan := make(chan os.Signal, 1)
	reload_timer := time.NewTicker(conf.PullInterval)
	reload_chan := make(chan bool)

	signal.Notify(signal_chan, os.Interrupt)
	defer reload_timer.Stop()

	engine := html.New("./views", ".html")
	engine.AddFunc("urljoin", util.UrlJoin)
	engine.AddFunc("timestr", conf.TimeStr)
	engine.AddFunc("host", util.Host)
	engine.AddFunc("first", util.First)
	engine.AddFunc("html", util.Html)
	engine.AddFunc("join", path.Join)
	engine.AddFunc("base", path.Base)
	engine.AddFunc("dir", path.Dir)
	engine.AddFunc("l", loc.Get)

	app = fiber.New(fiber.Config{
		AppName:               "teawiki",
		Views:                 engine,
		ServerHeader:          "",
		DisableStartupMessage: true,
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("reload", reload_chan)
		c.Locals("config", conf)
		c.Locals("locale", loc)
		c.Locals("repo", rep)
		return c.Next()
	})
	app.Get("/robots.txt", func(c *fiber.Ctx) error {
		return c.SendFile("./static/robots.txt")
	})
	app.Static("/_", "./static")

	// routes
	app.Get("/", routes.GET_Index)
	app.Post("/_/search", routes.POST_search)
	app.Post("/_/webhook", routes.POST_Webhook)
	app.Get("/_/history/*", routes.GET_History)
	app.Get("/*", routes.GET_page)

	// background thread
	go func() {
		var err error

		for {
			select {
			case <-reload_chan:
			case <-reload_timer.C:

			case <-signal_chan:
				goto shutdown
			}

			if err = rep.Reload(); err != nil {
				log.Fail("failed to reload the repo: %s", err.Error())
				break
			}
		}

	shutdown:
		reload_timer.Stop()
		close(signal_chan)
		close(reload_chan)

		log.Info("stopping the HTTP server")
		_ = app.Shutdown()
		os.Exit(1)
	}()

	// start the server
	log.Info("starting the HTTP server on %s", conf.ListenAddr)

	if err = app.Listen(conf.ListenAddr); err != nil {
		log.Fail("failed to start the HTTP server: %s", err.Error())
		os.Exit(1)
	}
}
