package models

import "gitlab.com/jfaucherre/mergo/git"

type PrContent interface {
	GetHead() string
	GetBase() string
	GetHost() string

	GetRepository() *git.Repo
	GetRemote() string
	GetRepoName() string
	GetOwner() string
	GetClipboard() bool
}

type Host interface {
	SubmitPr(PrContent) (*MRInfo, error)
}

type MRInfo struct {
	URL string
}
