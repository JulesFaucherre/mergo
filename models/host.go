package models

type MRInfo struct {
	URL string
}

// Host is the interface for hosts to be used
type Host interface {
	SubmitPr(*Opts) (*MRInfo, error)
}
