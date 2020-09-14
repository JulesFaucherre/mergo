package git

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"gitlab.com/jfaucherre/mergo/logger"
)

type Repo struct {
	path string
	repo *git.Repository
	wt   *git.Worktree
}

func New(p string) (*Repo, error) {
	p, err := getGitPath(p)
	if err != nil {
		return nil, err
	}

	r, err := git.PlainOpen(p)
	if err != nil {
		return nil, err
	}

	wt, err := r.Worktree()
	if err != nil {
		return nil, err
	}

	return &Repo{repo: r, wt: wt, path: p}, nil
}

func FromPwd() (*Repo, error) {
	p, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return New(p)
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
		return "", ErrGitRepositoryNotFound
	}

	return getGitPath(path.Join(repoPath, ".."))
}

func (me *Repo) HasChanges() bool {
	status, _ := me.wt.Status()
	logger.Debug("changes:\n%+v\n", status)
	return len(status) != 0
}

func (me *Repo) GetRemoteURLs(remote string) ([]string, error) {
	rmt, err := me.repo.Remote(remote)
	if err == git.ErrRemoteNotFound {
		return nil, ErrRemoteNotFound
	}
	if err != nil {
		return nil, err
	}

	return rmt.Config().URLs, nil
}

func (me *Repo) GetActualBranch() (string, error) {
	hd, err := me.repo.Head()
	if err != nil {
		return "", err
	}
	refName := string(hd.Name())
	logger.Debug("head reference %s\n", refName)
	if !strings.HasPrefix(refName, "refs/heads/") {
		return "", ErrNotABranch
	}
	// We remove the "refs/heads/" part
	name := refName[11:]
	return name, nil
}

func (me *Repo) GetBranchCommit(branch string) (*object.Commit, error) {
	ref, err := me.getBranchReference(branch)
	if err != nil {
		return nil, err
	}

	return me.repo.CommitObject(ref.Hash())
}

func (me *Repo) getBranchReference(branch string) (*plumbing.Reference, error) {
	var err error
	var b *plumbing.Reference

	branches, err := me.repo.Branches()
	if err != nil {
		return nil, err
	}
	for b, err = branches.Next(); b != nil && err == nil; b, err = branches.Next() {
		if strings.HasSuffix(string(b.Name()), branch) {
			return b, nil
		}
	}
	if err == io.EOF {
		return nil, ErrBranchNotFound
	}
	return nil, err
}

func (me *Repo) GetDifferenceCommit(hd, bs string) ([]string, error) {
	return nil, nil
}

func (me *Repo) IsBranchUpToDate(branch string) (bool, error) {
	rmts, err := me.repo.Remotes()
	if err != nil {
		return false, nil
	}
	for _, rmt := range rmts {
		up, err := me.isBranchUpToDate(branch, rmt.Config().Name)
		if err != nil {
			return false, err
		}
		if !up {
			return false, nil
		}
	}

	return true, nil
}

func (me *Repo) isBranchUpToDate(branch, remote string) (bool, error) {
	rev := plumbing.Revision(remote + "/" + branch)
	rmtH, err := me.repo.ResolveRevision(rev)
	if err != nil {
		return false, err
	}

	lclRef, err := me.getBranchReference(branch)
	if err != nil {
		return false, err
	}

	lclH := lclRef.Hash()
	return lclH == *rmtH, nil
}

func IsDirectChild(parent, child *object.Commit) bool {
	parents := child.Parents()

	for p, err := parents.Next(); p != nil && err == nil; p, err = parents.Next() {
		if p.Hash == parent.Hash {
			return true
		}
	}
	return false
}
