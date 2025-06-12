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

func matchesTitle(page *repo.Page, term string, exact bool) bool {
	if exact {
		// check if exact match
		return page.Title == term
	}

	// check if the title contains the term
	return strings.Contains(strings.ToLower(page.Title), term)
}

func matchesHeadings(headings []*repo.Heading, term string, exact bool) string {
	for _, heading := range headings {
		if exact {
			// check if exact match
			if heading.Name == term {
				return heading.ID
			}

			continue
		}

		// check if the heading contains the term
		if strings.Contains(strings.ToLower(heading.Name), term) {
			return heading.ID
		}

		// check the children headings
		if id := matchesHeadings(heading.Children, term, exact); id != "" {
			return id
		}
	}

	// no matches
	return ""
}

func matchesTags(page *repo.Page, term string, exact bool) bool {
	for _, tag := range page.Tags {
		if exact {
			// check if exact match
			if tag == term {
				return true
			}

			continue
		}

		// check if tag contains the term
		lower := strings.ToLower(tag)
		return strings.Contains(lower, term)
	}

	return false
}

func isType(list *map[string]*repo.Page, a, b int) bool {
	if a == b || a == SEARCH_ALL {
		*list = make(map[string]*repo.Page)
		return true
	}

	return false
}

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
			search_term = strings.TrimPrefix(term, "title:")
		case "heading":
			search_type = SEARCH_HEADING
			search_term = strings.TrimPrefix(term, "heading:")
		case "tag":
			search_type = SEARCH_TAG
			search_term = strings.TrimPrefix(term, "tag:")
		}
	}

	if !exact {
		search_term = strings.ToLower(search_term)
	}

	// search page titles
	if isType(&titles, search_type, SEARCH_TITLE) {
		rep.EachPage(func(p *repo.Page) {
			// check if the page title matches the search term
			if matchesTitle(p, search_term, exact) {
				titles[p.Relpath] = p
			}
		})
	}

	// search page headings
	if isType(&headings, search_type, SEARCH_HEADING) {
		rep.EachPage(func(p *repo.Page) {
			// check if any of the headings matches the term
			if id := matchesHeadings(p.Headings, search_term, exact); id != "" {
				headings[p.Relpath+"#"+id] = p
			}
		})
	}

	// search page tags
	if isType(&tags, search_type, SEARCH_TAG) {
		rep.EachPage(func(p *repo.Page) {
			// check if any of the page tags contain the term
			if matchesTags(p, search_term, exact) {
				tags[p.Relpath] = p
			}
		})
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
