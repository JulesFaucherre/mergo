package models

type Opts struct {
	Head string `short:"h" long:"head" description:"The head branch you want to merge into the base"`
	Base string `short:"b" long:"base" description:"The base branch you want to merge into" default:"master"`
	Host string `short:"g" long:"host" description:"The git host you use, ie github, gitlab, etc."`

	Remote string `short:"r" long:"remote" description:"The remote to use" default:"origin"`
	Repo   string `short:"d" long:"repository" description:"The name of the repository on which you want to make the pull request"`
	Owner  string `short:"w" long:"owner" description:"The owner of the repository"`

	Title string `short:"t" long:"title" description:"The title of the pull request"`
	Body  string `short:"c" long:"content" description:"The body of the pull request"`
}
