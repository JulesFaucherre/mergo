package models

type MRParams struct {
	URL     string
	Head    string
	Base    string
	Message string
}

type MRInfo struct {
	URL string
}

// Host is the interface for hosts to be used
type Host interface {
	SubmitPr(*MRParams) (*MRInfo, error)
}
