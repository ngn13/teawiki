package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/teawiki/repo"
	"github.com/ngn13/teawiki/util"
)

func GET_page(c *fiber.Ctx) error {
	rep := c.Locals("repo").(*repo.Repo)

	path := rep.Resolve(c.Path())
	page := rep.Get(path)

	if page == nil {
		return util.Send(c, path)
	}

	return util.Ok(c, "page", fiber.Map{
		"page": page,
		"path": path,
	})
}
