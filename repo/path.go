package repo

import (
	"path"
	"sort"
	"strings"

	"github.com/ngn13/teawiki/util"
)

func (r *Repo) Resolve(rp string) (string, bool) {
	rp = strings.ReplaceAll(rp, "..", "") // yeah nice try

	if rp[0] != '/' {
		rp = "/" + rp
	}

	base := path.Base(rp)

	if path.Ext(base) != "" {
		return rp, false
	}

	fp := path.Join(r.Config.RepoPath, rp)

	if util.IsDir(fp) {
		return path.Join(rp, INDEX_NAME), true
	}

	dir := path.Dir(rp)
	return path.Join(dir, base+PAGE_EXT), false
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
