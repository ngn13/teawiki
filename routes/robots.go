package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/teawiki/config"
)

func GET_Robots(c *fiber.Ctx) error {
	conf := c.Locals("config").(*config.Config)
	sitemap := conf.Url.JoinPath("sitemap.xml").String()

	robots := "User-agent: *\n"
	robots += "Disallow: /_/\n"
	robots += fmt.Sprintf("Sitemap: %s", sitemap)

	return c.SendString(robots)
}
