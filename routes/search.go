package routes

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/ngn13/teawiki/repo"
	"github.com/ngn13/teawiki/util"
)

const (
	SEARCH_TITLE   = 0 // search only page titles
	SEARCH_HEADING = 1 // search only page headings (h1, h2 etc.)
	SEARCH_TAG     = 2 // search only page tags
	SEARCH_ALL     = 3 // search all the stuff above
)

type headingResult struct {
	Heading *repo.Heading
	Page    *repo.Page
}

type matcher struct {
	term    string
	exact   bool
	current int
	results []map[string]interface{}
}

func (m *matcher) Set(i int) bool {
	if i >= len(m.results) {
		return false
	}

	m.current = i
	return true
}

func (m *matcher) Get(i int) map[string]interface{} {
	if i >= len(m.results) {
		return nil
	}

	return m.results[i]
}

func (m *matcher) Matches(target string, path string, res interface{}) bool {
	if m.exact && target != m.term {
		return false
	}

	if !m.exact && !strings.Contains(strings.ToLower(target), m.term) {
		return false
	}

	m.results[m.current][path] = res
	return true
}

func newMatcher(term string, exact bool, results int) *matcher {
	m := &matcher{
		current: 0,
		exact:   exact,
		results: make([]map[string]interface{}, results),
	}

	for i := range results {
		m.results[i] = make(map[string]interface{})
	}

	if m.exact {
		m.term = term
	} else {
		m.term = strings.ToLower(term)
	}

	return m
}

func findHeadings(m *matcher, p *repo.Page, hl []*repo.Heading) {
	headings := p.Headings

	if hl != nil {
		headings = hl
	}

	for _, h := range headings {
		// check if the heading matches
		m.Matches(h.Name, fmt.Sprintf("%s#%s", p.Relpath, h.ID), headingResult{
			Heading: h,
			Page:    p,
		})

		// check the children headings as well
		if h.Children != nil {
			findHeadings(m, p, h.Children)
		}
	}
}

func findTags(m *matcher, p *repo.Page) {
	for _, t := range p.Tags {
		// check if tags matches
		m.Matches(t, fmt.Sprintf("/_/tag/%s", t), t)
	}
}

func isType(a, b int) bool {
	return a == b || a == SEARCH_ALL
}

func search(c *fiber.Ctx, term string, exact bool) error {
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

	// create a new matcher
	matcher := newMatcher(search_term, exact, 3)

	if isType(search_type, SEARCH_TITLE) {
		// set the matcher result position
		matcher.Set(SEARCH_TITLE)

		// search & finding pages with matching titles
		rep.EachPage(func(p *repo.Page) {
			matcher.Matches(p.Title, p.Relpath, p)
		})
	}

	if isType(search_type, SEARCH_HEADING) {
		matcher.Set(SEARCH_HEADING)

		// search & finding matching headings
		rep.EachPage(func(p *repo.Page) {
			findHeadings(matcher, p, nil)
		})
	}

	if isType(search_type, SEARCH_TAG) {
		matcher.Set(SEARCH_TAG)

		// search & find matching tags
		rep.EachPage(func(p *repo.Page) {
			findTags(matcher, p)
		})
	}

	return util.Ok(c, "search", fiber.Map{
		"term": term,
		"type": search_type,
		"all":  search_type == SEARCH_ALL,
		"result": fiber.Map{
			"titles":   matcher.Get(SEARCH_TITLE),
			"headings": matcher.Get(SEARCH_HEADING),
			"tags":     matcher.Get(SEARCH_TAG),
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
