package models

type Opts struct {
	Head string `short:"d" long:"head" description:"The head branch you want to merge into the base"`
	Base string `short:"b" long:"base" description:"The base branch you want to merge into" default:"master"`
	Host string `long:"host" description:"The git host you use, ie github, gitlab, etc."`

	Remote string `long:"remote" description:"The remote to use" default:"origin"`
	Repo   string `long:"repository" description:"The name of the repository on which you want to make the pull request"`
	Owner  string `long:"owner" description:"The owner of the repository"`

	Delete string `long:"delete-creds" description:"Use this option when you mergo to delete the credentials it has stored about an host"`
}
