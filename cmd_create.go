package main

import (
	"context"
	"fmt"
	"regexp"

	"github.com/atotto/clipboard"
	"gitlab.com/jfaucherre/mergo/git"
	"gitlab.com/jfaucherre/mergo/hosts"
	"gitlab.com/jfaucherre/mergo/tools"
)

type CreateOptions struct {
	Head string `short:"d" long:"head" description:"The head branch you want to merge into the base"`
	Base string `short:"b" long:"base" description:"The base branch you want to merge into"`
	Host string `long:"host" description:"The git host you use, ie github, gitlab, etc."`

	Remote    string `long:"remote" description:"The remote to use"`
	Repo      string `long:"repository" description:"The name of the repository on which you want to make the pull request"`
	Owner     string `long:"owner" description:"The owner of the repository"`
	Clipboard bool   `long:"copy-clipboard" description:"Copies the merge request adress to the clipboard"`

	gitRepo *git.Repo
}

var (
	httpsR        = regexp.MustCompile(`https://(\w+\.\w+)/([\w-]+)/([\w-]+).git`)
	sshR          = regexp.MustCompile(`git@(\w+\.\w+):([\w-]+)/([\w-]+).git`)
	createOptions = &CreateOptions{
		Base:   "master",
		Remote: "origin",
	}
)

func init() {
	parser.AddCommand("create",
		"Create a pull request",
		"",
		createOptions)
}

func (me *CreateOptions) getGitRepo() {
	if me.gitRepo == nil {
		me.gitRepo = git.LocalRepository()
	}
}

func (me *CreateOptions) getRepoInfos() error {
	remoteString, err := me.
		gitRepo.
		Remote(me.GetRemote()).
		Do(context.Background())
	if err != nil {
		return err
	}

	matches := httpsR.FindStringSubmatch(remoteString)
	if matches == nil {
		matches = sshR.FindStringSubmatch(remoteString)
		if matches == nil {
			return fmt.Errorf("Unable to extract informations from remote string %s", remoteString)
		}
	}

	if tools.IsEmpty(me.Host) {
		me.Host = matches[1]
	}
	if tools.IsEmpty(me.Owner) {
		me.Owner = matches[2]
	}
	if tools.IsEmpty(me.Repo) {
		me.Repo = matches[3]
	}

	return nil
}

func (me *CreateOptions) GetHead() string {
	me.getGitRepo()
	if tools.IsEmpty(me.Head) {
		me.Head, _ = me.gitRepo.Branch().Do(context.Background())
	}
	return me.Head
}

func (me *CreateOptions) GetBase() string {
	me.getGitRepo()
	// TODO get default branch
	return me.Base
}

func (me *CreateOptions) GetHost() string {
	me.getGitRepo()
	me.getRepoInfos()
	return me.Host
}

func (me *CreateOptions) GetRemote() string {
	me.getGitRepo()
	return me.Remote
}

func (me *CreateOptions) GetRepoName() string {
	me.getGitRepo()
	me.getRepoInfos()
	return me.Repo
}

func (me *CreateOptions) GetOwner() string {
	me.getGitRepo()
	me.getRepoInfos()
	return me.Owner
}

func (me *CreateOptions) GetClipboard() bool {
	me.getGitRepo()
	return me.Clipboard
}

func (me *CreateOptions) GetRepository() *git.Repo {
	me.getGitRepo()
	return me.gitRepo
}

func (me *CreateOptions) Execute(args []string) error {
	fmt.Println("hello")
	host, err := hosts.GetHost(me.GetHost())
	fmt.Println("hello")
	if err != nil {
		return err
	}

	r, _ := host.SubmitPr(me)

	if me.GetClipboard() {
		clipboard.WriteAll(r.URL)
	}

	return nil
}
