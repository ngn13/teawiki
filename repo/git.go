package repo

import (
	"path"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
)

type History struct {
	Commit  string
	Message string
	Author  string
	Time    time.Time
}

func (r *Repo) pull() error {
	if r.Config.RepoUrl == nil {
		return nil
	}

	return r.Tree.Pull(&git.PullOptions{
		RemoteURL:  r.Config.RepoUrl.String(),
		RemoteName: "origin",
	})
}

func (r *Repo) history(file string, start, count int) ([]History, bool, error) {
	it, err := r.Git.Log(&git.LogOptions{
		PathFilter: func(p string) bool {
			return path.Join(r.Config.RepoPath, p) == file
		},
		From:  r.Head.Hash(),
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return nil, false, err
	}

	defer it.Close()

	history := []History{}
	more := false

	for commit, err := it.Next(); commit != nil && err == nil; commit, err = it.Next() {
		if start > 0 {
			start--
			continue
		}

		if count <= 0 {
			more = true
			break
		}

		history = append(history, History{
			Author:  commit.Author.String(),
			Message: strings.Split(commit.Message, "\n")[0],
			Commit:  commit.ID().String(),
			Time:    commit.Author.When,
		})
		count--
	}

	if err != nil {
		return nil, false, err
	}

	return history, more, nil
}
