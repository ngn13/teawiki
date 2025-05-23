package routes

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/teawiki/repo"
	"github.com/ngn13/teawiki/util"
)

const COMMIT_COUNT = 20

func GET_History(c *fiber.Ctx) error {
	n := c.QueryInt("n", 0)

	if n < 0 {
		n = 0
	}

	rep := c.Locals("repo").(*repo.Repo)
	rp := strings.Replace(c.Path(), "/_/history", "", 1)

	path, _ := rep.Resolve(rp)
	page := rep.Get(path)

	if page == nil {
		return util.NotFound(c)
	}

	history, more := rep.History(path, n*COMMIT_COUNT, COMMIT_COUNT)

	if len(history) == 0 {
		return util.NotFound(c)
	}

	return util.Ok(c, "history", fiber.Map{
		"page":    page,
		"path":    path,
		"history": history,
		"more":    more,
		"n":       n,
	})
}
