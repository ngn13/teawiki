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

func (r *Repo) loadPage(fp string, defaults ...string) (*Page, error) {
	var (
		reader *util.Reader = nil
		buffer *util.Buffer = nil

		file *os.File = nil
		page Page

		pos int64 = 0
		err error = nil
	)

	// open the file
	if file, err = os.Open(fp); err != nil {
		return nil, err
	}
	defer file.Close()

	// create the buffer
	buffer = util.NewBuffer(5)

	// check for the metadata section
	if err = buffer.From(file, 4); err != nil {
		log.Debg("failed to read start of metadata in %s: %s", fp, err.Error())
		goto markdown
	}

	if buffer.String() != "---\n" {
		log.Debg("%s is missing metadata start sign", fp)
		goto markdown
	}

	// read the rest of the metadata section
	pos = int64(buffer.Len())
	buffer.Clear()

	for {
		if err = buffer.From(file, 1); err != nil {
			log.Debg("cannot read byte from %s: %s", fp, err.Error())
			break
		}

		pos++

		if buffer.String() != "\n---\n" {
			continue
		}

		start := pos - int64(buffer.Len())

		if start <= 0 {
			log.Debg("missing metadata section in %s", fp)
			break
		}

		if reader, err = util.NewReader(file, 0, start); err != nil {
			log.Debg("failed to create reader for %s: %s", fp, err.Error())
			break
		}
		defer reader.Close()

		break
	}

	if reader == nil {
		log.Debg("%s is missing metada, no reader created", fp)
		goto markdown
	}

	// parse the YAML metadata & check if it's valid
	if err = yaml.NewDecoder(reader).Decode(&page); err != nil {
		return nil, err
	}

	for _, tag := range page.Tags {
		if tag == "" || strings.ContainsAny(tag, "\"!'^+%&/()=?*\\#,") {
			return nil, fmt.Errorf("bad tag name: %s", tag)
		}
	}

markdown:
	// if no metadata is read, seek to start of file to parse it all as markdown
	if reader == nil {
		if _, err = file.Seek(0, 0); err != nil {
			return nil, fmt.Errorf("failed to seek to start: %s", err.Error())
		}
	}

	// parse the markdown content & check if it's valid
	page.Content = string(r.Markdown.Render(file))

	if page.Content == "" {
		return nil, fmt.Errorf("empty page content")
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

	// use the filename as the title if none specified
	if page.Title == "" {
		page.Title = path.Base(fp)
	}

	// load headings from the parsed markdown
	page.Headings = Headings(page.Content)

	return &page, nil
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
