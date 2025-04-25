package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/teawiki/repo"
	"github.com/ngn13/teawiki/util"
)

func POST_Search(c *fiber.Ctx) error {
	body := struct {
		Term string `form:"term"`
	}{}

	if err := c.BodyParser(&body); err != nil {
		return util.BadRequest(c)
	}

	if body.Term == "" || body.Term == " " {
		return c.Redirect("/")
	}

	rep := c.Locals("repo").(*repo.Repo)
	results := rep.Find(body.Term)

	return util.Ok(c, "search", fiber.Map{
		"results": results,
		"count":   len(results),
	})
}
