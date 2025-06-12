package main

import (
	"os"
	"os/signal"
	"path"
	"syscall"
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
	reload_timer := time.NewTicker(conf.ReloadInterval)
	reload_chan := make(chan bool)

	signal.Notify(signal_chan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer reload_timer.Stop()

	// setup custom engine functions
	engine := html.New("./views", ".html")
	engine.AddFunc("sanitize", util.Sanitize)
	engine.AddFunc("urljoin", util.UrlJoin)
	engine.AddFunc("timestr", conf.TimeStr)
	engine.AddFunc("host", util.Host)
	engine.AddFunc("first", util.First)
	engine.AddFunc("html", util.Html)
	engine.AddFunc("join", path.Join)
	engine.AddFunc("base", path.Base)
	engine.AddFunc("dir", path.Dir)
	engine.AddFunc("map", util.Map)
	engine.AddFunc("add", util.Add)
	engine.AddFunc("l", loc.Get)

	// setup the web application
	app = fiber.New(fiber.Config{
		AppName:               "teawiki",
		Views:                 engine,
		ServerHeader:          "",
		DisableStartupMessage: true,
	})

	app.Static("/_", "./static")
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("reload", reload_chan)
		c.Locals("config", conf)
		c.Locals("locale", loc)
		c.Locals("repo", rep)
		return c.Next()
	})

	// HTTP routes
	app.Get("/", routes.GET_Index)
	app.Get("/robots.txt", routes.GET_Robots)
	app.Get("/sitemap.xml", routes.GET_Sitemap)
	app.Post("/_/search", routes.POST_Search)
	app.Get("/_/search", routes.GET_Search)
	app.Post("/_/webhook/:platform", routes.POST_Webhook)
	app.Get("/_/history/*", routes.GET_History)
	app.Get("/_/license", routes.GET_License)
	app.Get("/*", routes.GET_Page)

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
