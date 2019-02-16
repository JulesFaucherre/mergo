package models

type Host interface {
	SubmitPr(*Opts) error
	GetOwnerAndRepo(string) (string, string)
}
