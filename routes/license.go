package routes

import (
	"os"
	"path"

	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/teawiki/log"
	"github.com/ngn13/teawiki/util"
)

func GET_License(c *fiber.Ctx) error {
	var scripts []string = []string{}
	entries, err := os.ReadDir("./static/js")

	if err != nil {
		log.Fail("failed to read the javascript directory: %s", err.Error())
		return util.ServerError(c)
	}

	for _, entry := range entries {
		scripts = append(scripts, path.Join("/_/js", entry.Name()))
	}

	return util.Ok(c, "license", fiber.Map{
		"scripts": scripts,
	})
}
