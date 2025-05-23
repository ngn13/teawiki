package routes

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/teawiki/repo"
	"github.com/ngn13/teawiki/util"
)

func GET_Page(c *fiber.Ctx) error {
	rep := c.Locals("repo").(*repo.Repo)

	path, dir := rep.Resolve(c.Path())
	page := rep.Get(path)

	if dir && !strings.HasSuffix(c.OriginalURL(), "/") {
		return c.Redirect(c.OriginalURL() + "/")
	}

	if page == nil {
		return util.Send(c, path)
	}

	return util.Ok(c, "page", fiber.Map{
		"page": page,
		"path": path,
	})
}
