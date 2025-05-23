package repo

import (
	"fmt"
	"sort"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/ngn13/teawiki/config"
	"github.com/ngn13/teawiki/consts"
	"github.com/ngn13/teawiki/locale"
	"github.com/ngn13/teawiki/util"
)

const (
	LATEST_MAX = 10
	PAGE_EXT   = ".md"

	INDEX_NAME   = "README" + PAGE_EXT
	LICENSE_NAME = "LICENSE" + PAGE_EXT

	INDEX_PATH   = "/" + INDEX_NAME
	LICENSE_PATH = "/" + LICENSE_NAME
)

type Repo struct {
	Config   *config.Config
	Locale   *locale.Locale
	Markdown *util.Markdown

	Git  *git.Repository
	Head *plumbing.Reference
	Tree *git.Worktree

	Pages   map[string]*Page
	Latest  []string
	Index   *Page
	License *Page
}

func New(conf *config.Config, loc *locale.Locale) (*Repo, error) {
	var (
		repo = Repo{
			Config:   conf,
			Locale:   loc,
			Markdown: util.NewMd(conf.ChromaStyle),
		}
		err error
	)

	if util.Exists(conf.RepoPath) {
		if repo.Git, err = git.PlainOpen(conf.RepoPath); err != nil {
			repo.Git = nil
		}
	}

	if repo.Git == nil {
		if conf.RepoUrl == nil {
			return nil, fmt.Errorf("please specify a valid git repo URL or a path")
		}

		repo.Git, err = git.PlainClone(conf.RepoPath, false, &git.CloneOptions{
			URL: conf.RepoUrl.String(),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to clone remote repository: %s",
				err.Error(),
			)
		}
	}

	if repo.Tree, err = repo.Git.Worktree(); err != nil {
		return nil, fmt.Errorf("failed to get the git work tree: %s", err.Error())
	}

	if err = repo.Reload(); err != nil {
		return nil, fmt.Errorf("failed to load the repository: %s", err.Error())
	}

	return &repo, nil
}

func (r *Repo) Reload() error {
	var err error

	if err = r.pull(); err != nil {
		return err
	}

	// obtain the head
	if r.Head, err = r.Git.Head(); err != nil {
		return err
	}

	// clear all the pages
	r.Pages = make(map[string]*Page)
	r.Latest = []string{}
	r.License = nil
	r.Index = nil

	// traverse & load all the pages
	if err = r.traverse(); err != nil {
		return err
	}

	paths := []string{}

	for path := range r.Pages {
		paths = append(paths, path)
	}

	// sort the paths based on the pages (latest first)
	sort.Slice(paths, func(i, j int) bool {
		switch r.Pages[paths[i]].LastUpdate.Compare(r.Pages[paths[j]].LastUpdate) {
		case 0:
			return r.Pages[paths[j]].Title > r.Pages[paths[i]].Title

		case -1:
			return false
		}

		return true
	})

	// create the latest updated page list
	for i, path := range paths {
		if i >= LATEST_MAX {
			break
		}
		r.Latest = append(r.Latest, path)
	}

	r.License = r.Get(LICENSE_PATH)
	if r.Index = r.Get(INDEX_PATH); r.Index == nil {
		r.Index = r.newPage(
			r.Locale.Get("index.title"),
			r.Locale.Get("index.content", consts.DOCS),
		)
	}

	return nil
}
