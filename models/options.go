package models

// Opts represents the options passed from command line
type CreateOptions struct {
	Head string
	Base string
	Host string

	Remote    string
	Repo      string
	Owner     string
	Clipboard bool
}
