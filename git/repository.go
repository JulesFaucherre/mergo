package git

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"gitlab.com/jfaucherre/mergo/tools"
)

type Repo struct {
	p string
}

// Repository returns the used repository at path 'repoPath'
func Repository(repoPath string) *Repo {
	return &Repo{p: repoPath}
}

// LocalRepository returns the used repository at pwd
func LocalRepository() *Repo {
	pwd, _ := os.Getwd()
	return Repository(pwd)
}

// Branch returns a GitCmd to get the active branch of the repository
func (me *Repo) Branch() *GitCmd {
	return &GitCmd{
		repo: me,
		cmd: [][]string{
			{"git", "branch"},
			{"grep", "*"},
			{"awk", "{print $2}"},
		},
	}
}

// Remote returns the git url for the remote 'remote'
func (me *Repo) Remote(remote string) *GitCmd {
	return &GitCmd{
		repo: me,
		cmd: [][]string{
			{"git", "remote", "get-url", remote},
		},
	}
}

type GitCmd struct {
	repo *Repo
	cmd  [][]string
	next func(string, error) (string, error)
}

// Do runs the GitCmd with the context ctx and returns its result
func (me *GitCmd) Do(ctx context.Context) (string, error) {
	if tools.IsEmpty(me.repo.p) || me.cmd == nil {
		return "", fmt.Errorf("Can not launch any command, repository was not initiated well")
	}
	repoPath, err := getGitPath(me.repo.p)
	if err != nil {
		return "", err
	}

	pwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if err = os.Chdir(repoPath); err != nil {
		return "", err
	}

	if me.next == nil {
		me.next = func(a string, e error) (string, error) { return a, e }
	}

	res, err := me.next(run(ctx, me.cmd))
	if err != nil {
		return "", err
	}

	if err = os.Chdir(pwd); err != nil {
		return "", err
	}

	return res, nil
}

func getGitPath(repoPath string) (string, error) {
	repoPath, err := filepath.Abs(repoPath)
	if err != nil {
		return "", err
	}

	gitPath := path.Join(repoPath, ".git")
	stat, err := os.Stat(gitPath)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}
	if !os.IsNotExist(err) && stat.IsDir() {
		return repoPath, nil
	}
	if repoPath == "/" {
		return "", fmt.Errorf("Unable to find git repository")
	}

	return getGitPath(path.Join(repoPath, ".."))
}
