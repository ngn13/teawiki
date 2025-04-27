package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/teawiki/config"
)

func GET_Robots(c *fiber.Ctx) error {
	conf := c.Locals("config").(*config.Config)

	robots := "User-agent: *\n"
	robots += "Disallow: /_/\n"

	if conf.Url != nil {
		sitemap := conf.Url.JoinPath("sitemap.xml").String()
		robots += fmt.Sprintf("Sitemap: %s", sitemap)
	}

	return c.SendString(robots)
}
