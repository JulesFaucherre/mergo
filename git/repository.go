package git

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

type GitCmd struct {
	p   string
	cmd [][]string
}

func Repository(repoPath string) *GitCmd {
	return &GitCmd{p: repoPath}
}

func LocalRepository() *GitCmd {
	pwd, _ := os.Getwd()
	return Repository(pwd)
}

func (me *GitCmd) Branch() *GitCmd {
	me.cmd = [][]string{
		{"git", "branch"},
		{"grep", "*"},
		{"awk", "{print $2}"},
	}
	return me
}

func (me *GitCmd) Remote(remote string) *GitCmd {
	me.cmd = [][]string{
		{"git", "remote", "get-url", remote},
	}
	return me
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

func (me *GitCmd) Do(ctx context.Context) (string, error) {
	if me.p == "" || me.cmd == nil {
		return "", fmt.Errorf("Can not launch any command, repository was not initiated well")
	}
	_, err := getGitPath(me.p)
	if err != nil {
		return "", err
	}

	res, err := run(ctx, me.cmd)

	return res, err
}
