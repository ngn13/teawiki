package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/teawiki/repo"
	"github.com/ngn13/teawiki/util"
)

func GET_Page(c *fiber.Ctx) error {
	rep := c.Locals("repo").(*repo.Repo)

	cpath := c.Path()
	rpath := rep.Resolve(cpath)

	// check if we failed to resolve the path
	if rpath == "" {
		return util.NotFound(c)
	}

	// check if the requested patch matches with the resolved path
	if rpath != cpath {
		return c.Redirect(rpath)
	}

	// obtain the page from the resolved path
	page := rep.Get(rpath)

	if page == nil {
		return util.Send(c, rpath)
	}

	return util.Ok(c, "page", fiber.Map{
		"page": page,
		"path": rpath,
	})
}
