package routes

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/teawiki/repo"
	"github.com/ngn13/teawiki/util"
)

const (
	SEARCH_TITLE   = 0
	SEARCH_HEADING = 1
	SEARCH_TAG     = 2
	SEARCH_ALL     = 3
)

func search(c *fiber.Ctx, term string, exact bool) error {
	var titles, headings, tags map[string]*repo.Page

	rep := c.Locals("repo").(*repo.Repo)
	target := strings.Split(term, ":")

	search_type := SEARCH_ALL
	search_term := term

	if len(target) == 2 {
		// get the search type
		switch target[0] {
		case "title":
			search_type = SEARCH_TITLE
			search_term = strings.TrimLeft(term, "title:")
		case "heading":
			search_type = SEARCH_HEADING
			search_term = strings.TrimLeft(term, "heading:")
		case "tag":
			search_type = SEARCH_TAG
			search_term = strings.TrimLeft(term, "tag:")
		}
	}

	if search_type == SEARCH_TITLE || search_type == SEARCH_ALL {
		titles = rep.SearchTitles(search_term, exact)
	}

	if search_type == SEARCH_HEADING || search_type == SEARCH_ALL {
		headings = rep.SearchHeadings(search_term, exact)
	}

	if search_type == SEARCH_TAG || search_type == SEARCH_ALL {
		tags = rep.SearchTags(search_term, exact)
	}

	return util.Ok(c, "search", fiber.Map{
		"term": term,
		"result": fiber.Map{
			"titles":   titles,
			"headings": headings,
			"tags":     tags,
		},
		"length": fiber.Map{
			"titles":   len(titles),
			"headings": len(headings),
			"tags":     len(tags),
		},
	})
}

func POST_Search(c *fiber.Ctx) error {
	body := struct {
		Term  string `form:"term"`
		Exact bool   `form:"exact"`
	}{}

	if err := c.BodyParser(&body); err != nil {
		return util.BadRequest(c)
	}

	if body.Term == "" || body.Term == " " {
		return c.Redirect("/")
	}

	return search(c, body.Term, body.Exact)
}

func GET_Search(c *fiber.Ctx) error {
	term := c.Query("term")
	_exact := c.Query("exact")
	exact := false

	if term == "" {
		return c.Redirect("/")
	}

	if _exact == "1" || _exact == "true" {
		exact = true
	}

	return search(c, term, exact)
}
