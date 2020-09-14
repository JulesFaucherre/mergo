package git

import "errors"

var (
	ErrNotABranch            = errors.New("you must be on a branch to use it as default head")
	ErrBranchNotFound        = errors.New("branch not found")
	ErrGitRepositoryNotFound = errors.New("repository not found")
	ErrRemoteNotFound        = errors.New("remote not found")
)
