package routes

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/teawiki/repo"
	"github.com/ngn13/teawiki/util"
)

const COMMIT_COUNT = 20

func GET_History(c *fiber.Ctx) error {
	rep := c.Locals("repo").(*repo.Repo)
	n := c.QueryInt("n", 0)

	if n < 0 {
		n = 0
	}

	cpath := strings.TrimPrefix(c.Path(), "/_/history")
	rpath := rep.Resolve(cpath)

	// check if we failed to resolve the path
	if rpath == "" {
		return util.NotFound(c)
	}

	// check if the request path matches to resolved patch
	if rpath != cpath {
		return c.Redirect("/_/history" + rpath)
	}

	// obtain the page from the resolved path
	page := rep.Get(rpath)

	if page == nil {
		return util.NotFound(c)
	}

	// get the commit history of the page
	history, more := rep.History(rpath, n*COMMIT_COUNT, COMMIT_COUNT)

	if len(history) == 0 {
		return util.NotFound(c)
	}

	// TODO: maybe add a RSS feed to the history page?

	return util.Ok(c, "history", fiber.Map{
		"page":    page,
		"path":    rpath,
		"history": history,
		"more":    more,
		"n":       n,
	})
}
