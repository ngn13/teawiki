package routes

import (
	"encoding/xml"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/teawiki/config"
	"github.com/ngn13/teawiki/log"
	"github.com/ngn13/teawiki/util"
)

var osd_str string = ""

const (
	OPENSEARCH_URL = "http://a9.com/-/spec/opensearch/1.1/"
	INPUT_ENCODING = "UTF-8"
)

type opensearch_image struct {
	Url    string `xml:",chardata"`
	Width  uint   `xml:"width,attr"`
	Height uint   `xml:"height,attr"`
}

type opensearch_url struct {
	Type     string `xml:"type,attr"`
	Method   string `xml:"method,attr"`
	Template string `xml:"template,attr"`
}

type opensearch struct {
	XMLName       xml.Name `xml:"OpenSearchDescription"`
	Xmlns         string   `xml:"xmlns,attr"`
	ShortName     string
	Description   string
	InputEncoding string
	Image         *opensearch_image `xml:"Image,omitempty"`
	Url           *opensearch_url
}

func GET_Opensearch(c *fiber.Ctx) error {
	conf := c.Locals("config").(*config.Config)

	// URL needs to be configured for this route
	if conf.Url == nil {
		return util.NotFound(c)
	}

	// this XML is static, no need to rebuild it every time, save it once built
	if osd_str == "" {
		search_url := util.UrlJoin(conf.Url, "/_/search")
		template_url := fmt.Sprintf("%s?term={searchTerms}", search_url)

		osd := opensearch{
			Xmlns:         OPENSEARCH_URL,
			ShortName:     conf.Name,
			Description:   conf.Desc,
			InputEncoding: INPUT_ENCODING,
			Url: &opensearch_url{
				Type:     "text/html",
				Method:   "GET",
				Template: template_url,
			},
		}

		// add the icon if configured
		if conf.Icon != "" {
			osd.Image = &opensearch_image{
				Url:    util.UrlJoin(conf.Url, fmt.Sprintf("/_/assets/%s", conf.Icon)),
				Width:  16,
				Height: 16,
			}
		}

		body, err := xml.MarshalIndent(osd, "", "  ")
		if err != nil {
			log.Fail("failed to encode sitemap: %s", err.Error())
			return util.ServerError(c)
		}

		osd_str = string(body)
	}

	// send the XML contents
	c.Set("Content-Type", "text/xml; charset=utf-8")
	return c.SendString(xml.Header + osd_str)
}
