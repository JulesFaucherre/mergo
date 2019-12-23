package models

import "gitlab.com/jfaucherre/mergo/git"

// Opts represents the options passed from command line
type Opts struct {
	Head string `ini:"head" short:"d" long:"head" description:"The head branch you want to merge into the base"`
	Base string `ini:"base" short:"b" long:"base" description:"The base branch you want to merge into" default:"master"`
	Host string `ini:"host" long:"host" description:"The git host you use, ie github, gitlab, etc."`

	Remote     string `ini:"remote" long:"remote" description:"The remote to use" default:"origin"`
	Repository string `ini:"repository" long:"repository" description:"The name of the repository on which you want to make the pull request"`
	Owner      string `ini:"owner" long:"owner" description:"The owner of the repository"`
	Clipboard  bool   `ini:"clipboard" long:"copy-clipboard" description:"Copies the merge request adress to the clipboard"`

	Verbose bool `ini:"verbose" short:"v" long:"verbose" description:"Whether you want to have logs on whaat is happening"`

	Delete string `long:"delete-creds" description:"Use this option when you mergo to delete the credentials it has stored about an host"`

	Repo *git.Repo
}
