package routes

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/teawiki/repo"
	"github.com/ngn13/teawiki/util"
)

const (
	SEARCH_ALL     = 0
	SEARCH_TITLE   = 1
	SEARCH_HEADING = 2
	SEARCH_TAG     = 3
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

func matchesTags(page *repo.Page, term string, exact bool) string {
	for _, tag := range page.Tags {
		if exact {
			// check if exact match
			if tag == term {
				return tag
			}

			continue
		}

		// check if tag contains the term
		if strings.Contains(strings.ToLower(tag), term) {
			return tag
		}

		return ""
	}

	return ""
}

func isType(a, b int) bool {
	return a == b || a == SEARCH_ALL
}

func search(c *fiber.Ctx, term string, exact bool) error {
	var (
		titles, headings map[string]*repo.Page
		tags             map[string]string
	)

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
	if isType(search_type, SEARCH_TITLE) {
		// create the page URL to page map
		titles = make(map[string]*repo.Page)

		rep.EachPage(func(p *repo.Page) {
			// check if the page title matches the search term
			if matchesTitle(p, search_term, exact) {
				titles[p.Relpath] = p
			}
		})
	}

	// search page headings
	if isType(search_type, SEARCH_HEADING) {
		// create the heading URL to page map
		headings = make(map[string]*repo.Page)

		rep.EachPage(func(p *repo.Page) {
			// check if any of the headings matches the term
			if id := matchesHeadings(p.Headings, search_term, exact); id != "" {
				headings[p.Relpath+"#"+id] = p
			}
		})
	}

	// search page tags
	if isType(search_type, SEARCH_TAG) {
		// create the tag URL to tag
		tags = make(map[string]string)

		rep.EachPage(func(p *repo.Page) {
			// check if any of the page tags contain the term
			if tag := matchesTags(p, search_term, exact); tag != "" {
				tags["/_/tag/"+tag] = tag
			}
		})
	}

	return util.Ok(c, "search", fiber.Map{
		"term": term,
		"type": search_type,
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
