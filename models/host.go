package models

import "gitlab.com/jfaucherre/mergo/git"

type MRInfo struct {
	URL string
}

// Host is the interface for hosts to be used
type Host interface {
	SubmitPr(*CreateOptions, *git.Repo) (*MRInfo, error)
}
