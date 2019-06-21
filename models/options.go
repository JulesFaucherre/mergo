package models

type Options struct {
	Verbose []bool `short:"v" long:"verbose" description:"The level of verbosity you want"`
}

type ConfigOptions struct {
	Get        []string          `long:"get" description:"Gets a variable"`
	Set        map[string]string `long:"set" description:"Add a new variable like key:value"`
	Unset      []string          `long:"unset" description:"Unsets a variable"`
	DeleteCred []string          `long:"delete-credential" description:"Deletes the authentification credentials for the specified host"`
	Global     bool              `short:"g" long:"global" description:"Use global config file"`
}

type MergeOptions struct {
	Head string `short:"d" long:"head" description:"The head branch you want to merge into the base"`
	Base string `short:"b" long:"base" description:"The base branch you want to merge into"`
	Host string `long:"host" description:"The git host you use, ie github, gitlab, etc."`

	Remote    string `long:"remote" description:"The remote to use"`
	Repo      string `long:"repository" description:"The name of the repository on which you want to make the pull request"`
	Owner     string `long:"owner" description:"The owner of the repository"`
	Clipboard bool   `long:"copy-clipboard" description:"Copies the merge request adress to the clipboard"`
}

var (
	options       = &Options{Verbose: []bool{}}
	configOptions = &ConfigOptions{}
	mergeOptions  = &MergeOptions{
		Base:   "master",
		Remote: "origin",
	}
)
