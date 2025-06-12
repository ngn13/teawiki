package repo

import "strings"

func matchesTitle(page *Page, term string, exact bool) bool {
	if exact {
		// check if exact match
		return page.Title == term
	}

	// check if the title contains the term
	return strings.Contains(strings.ToLower(page.Title), term)
}

// search page titles and the paths for the provided term
func (r *Repo) SearchTitles(term string, exact bool) map[string]*Page {
	results := make(map[string]*Page)

	if !exact {
		term = strings.ToLower(term)
	}

	r.EachPage(func(p *Page) {
		// check if the page title matches the search term
		if matchesTitle(p, term, exact) {
			results[p.Relpath] = p
		}
	})

	return results
}

func matchesHeadings(headings []*Heading, term string, exact bool) string {
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

// search page headings for the provided term
func (r *Repo) SearchHeadings(term string, exact bool) map[string]*Page {
	results := make(map[string]*Page)

	if !exact {
		term = strings.ToLower(term)
	}

	r.EachPage(func(p *Page) {
		// check if any of the headings matches the term
		if id := matchesHeadings(p.Headings, term, exact); id != "" {
			results[p.Relpath+"#"+id] = p
		}
	})

	// return results
	return results
}

func matchesTags(page *Page, term string, exact bool) bool {
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

// search page tags for the provided term
func (r *Repo) SearchTags(term string, exact bool) map[string]*Page {
	results := make(map[string]*Page)

	if !exact {
		term = strings.ToLower(term)
	}

	for _, page := range r.Pages {
		// check if any of the page tags contain the term
		if matchesTags(page, term, exact) {
			results[page.Relpath] = page
		}
	}

	// return results
	return results
}
