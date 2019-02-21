package models

// Host is the interface for hosts to be used
type Host interface {
	SubmitPr(*Opts) error
}
