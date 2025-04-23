package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/teawiki/repo"
	"github.com/ngn13/teawiki/util"
)

func GET_Index(c *fiber.Ctx) error {
	return util.Ok(c, "page", fiber.Map{
		"page": c.Locals("repo").(*repo.Repo).Index,
		"path": repo.INDEX_PATH,
	})
}
