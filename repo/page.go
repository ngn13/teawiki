package repo

import (
	"fmt"
	"os"
	"path"
	"slices"
	"strings"
	"time"

	"github.com/goccy/go-yaml"
	"github.com/ngn13/teawiki/log"
	"github.com/ngn13/teawiki/util"
)

type Field struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
	Link  string `yaml:"link"`
}

type Page struct {
	Relpath    string     `yaml:"-"`
	Abspath    string     `yaml:"-"`
	HasHistory bool       `yaml:"-"`
	LastUpdate time.Time  `yaml:"-"`
	Headings   []*Heading `yaml:"-"`

	// this stuff is parsed from the file
	Title   string   `yaml:"title"`
	Image   string   `yaml:"image"`
	Tags    []string `yaml:"tags"`
	Fields  []Field  `yaml:"fields"`
	Content string   `yaml:"-"`
}

func (p *Page) IsValid() bool {
	return p.Title != "" && p.Content != ""
}

func (p *Page) Path(id ...string) string {
	if len(id) <= 0 {
		return p.Relpath
	}

	return fmt.Sprintf("%s#%s", p.Relpath, id[0])
}

func (r *Repo) loadPage(fp string, defaults ...string) (page *Page, err error) {
	var (
		yaml_reader *util.Reader
		mark_reader *util.Reader
		file        *os.File
	)

	if file, err = os.Open(fp); err != nil {
		return nil, err
	}
	defer file.Close()

	buff := util.NewBuffer(5)
	char := make([]byte, 1)
	pos := int64(0)

	for _, err = file.Read(char); err == nil; _, err = file.Read(char) {
		pos++

		if buff.Push(char[0]) != buff.Length() {
			continue
		}

		if buff.String() != "\n%%%\n" {
			continue
		}

		start := pos - int64(buff.Length())

		if start <= 0 {
			return nil, fmt.Errorf("invalid page format")
		}

		yaml_file, _ := os.Open(fp)
		mark_file, _ := os.Open(fp)

		if yaml_reader, err = util.NewReader(yaml_file, 0, start); err != nil {
			return nil, err
		}
		defer yaml_reader.Close()

		if mark_reader, err = util.NewReader(mark_file, pos, 0); err != nil {
			return nil, err
		}
		defer mark_reader.Close()

		break
	}

	if mark_reader == nil || yaml_reader == nil {
		return nil, fmt.Errorf("failed to parse the page")
	}

	page = &Page{}

	// parse the YAML metadata
	if err = yaml.NewDecoder(yaml_reader).Decode(page); err != nil {
		return nil, err
	}

	// parse the markdown content
	page.Content = string(r.Markdown.Render(mark_reader))

	if !page.IsValid() {
		return nil, fmt.Errorf("invalid page data")
	}

	// get the last update time
	if history, _, err := r.history(fp, 0, 1); err != nil {
		return nil, err
	} else if len(history) > 0 {
		page.HasHistory = true
		page.LastUpdate = history[0].Time
	} else {
		page.HasHistory = false
		page.LastUpdate = time.Unix(0, 0)
	}

	// load headings from the parsed markdown
	page.Headings = Headings(page.Content)

	return page, nil
}

func (r *Repo) newPage(title string, content string) *Page {
	html := string(r.Markdown.Render(strings.NewReader(content)))

	return &Page{
		HasHistory: false,
		LastUpdate: time.Unix(0, 0),
		Headings:   Headings(html),
		Title:      title,
		Content:    html,
	}
}

// traverses the git repo recursively to obtain all the pages
func (r *Repo) traverse(dir ...string) error {
	var (
		page    *Page
		entries []os.DirEntry
		err     error
	)

	rel_dir := "/"

	if len(dir) >= 1 {
		rel_dir = dir[0]
	}

	abs_dir := path.Join(r.Config.RepoPath, rel_dir)

	if entries, err = os.ReadDir(abs_dir); err != nil {
		return err
	}

	for _, entry := range entries {
		name := entry.Name()

		if name[0] == '.' || name[0] == '_' {
			continue
		}

		rel_entry := path.Join(rel_dir, name)
		abs_entry := path.Join(abs_dir, name)
		page = nil

		if entry.IsDir() {
			// if entry is a directory, traverse it
			err = r.traverse(rel_entry)
		} else if path.Ext(name) == PAGE_EXT {
			// otherwise just load the page
			page, err = r.loadPage(abs_entry)
		}

		// check if failed to load the page
		if err != nil {
			log.Warn("failed to load %s: %s", rel_entry, err.Error())
			continue
		}

		// successfuly loaded the page, save it
		if page != nil {
			page.Relpath = rel_entry
			page.Abspath = abs_entry
			r.Pages = append(r.Pages, page)
		}
	}

	return nil
}

func (r *Repo) Get(rp string) *Page {
	indx := slices.IndexFunc(r.Pages, func(p *Page) bool {
		return p.Relpath == rp
	})

	if indx < 0 {
		return nil
	}

	return r.Pages[indx]
}

func (r *Repo) History(fp string, start, count int) ([]History, bool) {
	rfp := path.Join(r.Config.RepoPath, fp)

	if history, more, err := r.history(rfp, start, count); err != nil {
		log.Fail("failed to get history for %s: %s", fp, err.Error())
		return nil, false
	} else {
		return history, more
	}
}
