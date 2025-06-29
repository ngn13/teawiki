package routes

import (
	"encoding/xml"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/teawiki/config"
	"github.com/ngn13/teawiki/log"
	"github.com/ngn13/teawiki/repo"
	"github.com/ngn13/teawiki/util"
)

const (
	W3C_DATETIME = "2006-01-02T15:04:05-07:00"
	SITEMAP_URL  = "http://www.sitemaps.org/schemas/sitemap/0.9"
)

type sitemap_url struct {
	XMLName  xml.Name `xml:"url"`
	Location string   `xml:"loc"`
	LastMod  string   `xml:"lastmod,omitempty"`
	Priority string   `xml:"priority,omitempty"`
}

type sitemap struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	Urls    []sitemap_url
}

func (s *sitemap) Add(loc string, mod *time.Time, _prio ...string) {
	prio := ""
	lastmod := ""

	if len(_prio) > 0 {
		prio = _prio[0]
	}

	if mod != nil {
		lastmod = mod.Format(W3C_DATETIME)
	}

	s.Urls = append(s.Urls, sitemap_url{
		Location: loc,
		LastMod:  lastmod,
		Priority: prio,
	})
}

func GET_Sitemap(c *fiber.Ctx) error {
	conf := c.Locals("config").(*config.Config)
	rep := c.Locals("repo").(*repo.Repo)

	if conf.Url == nil {
		return util.NotFound(c)
	}

	set := sitemap{
		Xmlns: SITEMAP_URL,
		Urls:  []sitemap_url{},
	}

	if rep.Index.HasHistory {
		set.Add(conf.Url.String(), &rep.Index.LastUpdate, "1.0")
	} else {
		set.Add(conf.Url.String(), nil, "1.0")
	}

	for _, page := range rep.Pages {
		if !page.HasHistory || page.Relpath == repo.INDEX_PATH {
			continue
		}

		fp := conf.Url.JoinPath(page.Relpath).String()
		set.Add(fp, &page.LastUpdate)
	}

	body, err := xml.MarshalIndent(set, "", "  ")
	if err != nil {
		log.Fail("failed to encode sitemap: %s", err.Error())
		return util.ServerError(c)
	}

	c.Set("Content-Type", "text/xml; charset=utf-8")
	return c.SendString(xml.Header + string(body))
}
