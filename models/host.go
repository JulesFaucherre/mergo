package models

type Host interface {
	SubmitPr(*Opts) error
}
