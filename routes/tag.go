package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/teawiki/repo"
	"github.com/ngn13/teawiki/util"
)

func GET_Tag(c *fiber.Ctx) error {
	rep := c.Locals("repo").(*repo.Repo)
	tag := c.Params("tag")
	pages := []*repo.Page{}

	// find pages with the tag
	rep.EachPage(func(p *repo.Page) {
		for _, t := range p.Tags {
			if t == tag {
				pages = append(pages, p)
			}
		}
	})

	// check if we found any pages
	if len(pages) == 0 {
		return util.NotFound(c)
	}

	return util.Ok(c, "tag", fiber.Map{
		"tag":   tag,
		"pages": pages,
	})
}
