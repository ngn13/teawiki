package repo

import (
	"path"
	"sort"
	"strings"

	"github.com/ngn13/teawiki/util"
)

func (r *Repo) Resolve(rp string) string {
	rp = strings.ReplaceAll(rp, "..", "") // yeah nice try

	if rp[0] != '/' {
		rp = "/" + rp
	}

	base := path.Base(rp)

	if path.Ext(base) == PAGE_EXT {
		return rp
	}

	fp := path.Join(r.Config.RepoPath, rp)

	if util.IsDir(fp) {
		return path.Join(rp, INDEX_NAME)
	}

	dir := path.Dir(rp)
	return path.Join(dir, base+PAGE_EXT)
}

func (r *Repo) List(dir string) []string {
	results := []string{}

	for pth := range r.Pages {
		d := path.Dir(pth)

		if d == dir {
			results = append(results, pth)
			continue
		}

		if path.Base(pth) == INDEX_NAME && path.Dir(d) == dir {
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
