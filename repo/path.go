package repo

import (
	"path"
	"sort"
	"strings"

	"github.com/ngn13/teawiki/util"
)

func (r *Repo) Resolve(rp string) string {
	rp = strings.ReplaceAll(rp, "..", "") // yeah nice try

	// path should always start with "/"
	if rp[0] != '/' {
		rp = "/" + rp
	}

	// check the root dir name, "_" is not allowed in root dir
	names := strings.Split(rp, "/")

	for _, name := range names {
		if name == "" {
			continue // skip empty names
		}

		if name == "_" {
			return "" // invalid path
		}

		break
	}

	// get the base name (file name)
	base := path.Base(rp)

	if path.Ext(base) != "" {
		return rp
	}

	// get full path of page and check if it's a directory
	fp := path.Join(r.Config.RepoPath, rp)

	if util.IsDir(fp) {
		return path.Join(rp, INDEX_NAME)
	}

	// if this is a file, add missing page extension
	dir := path.Dir(rp)
	return path.Join(dir, base+PAGE_EXT)
}

func (r *Repo) List(dir string) []string {
	results := []string{}

	for _, page := range r.Pages {
		d := path.Dir(page.Relpath)

		if d == dir {
			results = append(results, page.Relpath)
			continue
		}

		if path.Base(page.Relpath) == INDEX_NAME && path.Dir(d) == dir {
			results = append(results, d)
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i] == INDEX_PATH {
			return true
		}

		if results[j] == INDEX_PATH {
			return false
		}

		return results[i] < results[j]
	})

	return results
}
