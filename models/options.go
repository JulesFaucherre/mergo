package models

// Opts represents the options passed from command line
type Opts struct {
	Verbose []bool `short:"v" long:"verbose" description:"Add logs, you can have more logs by calling it more times" ini:"verbose"`

	Head string `short:"d" long:"head" description:"The head branch you want to merge into the base" default-mask:"the actual checked out branch" ini:"head" ini-name:"head"`
	Base string `short:"b" long:"base" description:"The base branch you want to merge into" default-mask:"master" ini:"base"`

	Message string `short:"m" long:"message" default-mask:"If you have only one commit, it takes this commit's message" description:"The pull request message"`

	Force bool `long:"force" short:"f" description:"Force the pull request, doesn't ask you if you have unstaged changes or things like that" ini:"force"`

	Clipboard bool `long:"clipboard" short:"c" description:"Copy the URLs of your merge requests to your clipboard" ini:"clipboard"`

	Remote     string   `long:"remote" description:"The remote to use" default-mask:"origin" ini:"remote"`
	RemoteURLs []string `short:"r" long:"remote-url" description:"The remote URLs to use. Note that this overwrite the \"remote\" option" ini:"remote-urls"`

	ForceEdit bool `short:"e" long:"force-edition" description:"Force the edition of the message event it already have a value" ini:"force-edition"`

	Delete string `long:"delete-creds" description:"Use this option when you want to delete the credentials of an host"`
}

// We are parsing ini files to have default values and then parse the options
// from the command line but the default values of go-flags overwrites the
// values from the ini configuration files so this function sets the default
// manually only when the fields are empty
func SetDefaultOptions(opts *Opts) {
	if opts.Remote == "" {
		opts.Remote = "origin"
	}
	if opts.Base == "" {
		opts.Base = "master"
	}
}
