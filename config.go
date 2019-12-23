package main

import (
	"fmt"
	"os"
	"path"

	flags "github.com/jessevdk/go-flags"
	"gitlab.com/jfaucherre/mergo/git"
	"gitlab.com/jfaucherre/mergo/models"
	"gitlab.com/jfaucherre/mergo/tools"
	ini "gopkg.in/ini.v1"
)

var (
	sectionName = "mergo"
)

func isSectionExists(e error) bool {
	return e.Error() == "section '"+sectionName+"' does not exist"
}

func loadFile(fname string, opts *models.Opts) error {
	// If there is no file, there is nothing to load
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		return nil
	}

	// Load ini file
	f, err := ini.Load(fname)
	if err != nil {
		return err
	}

	// Load section, if section does not exist return
	s, err := f.GetSection(sectionName)
	if err != nil && isSectionExists(err) {
		return nil
	}
	if err != nil {
		return err
	}

	return s.MapTo(opts)
}

func loadConfig() (*models.Opts, error) {
	opts := &models.Opts{}

	// Load global config
	home, _ := os.UserHomeDir()
	globalConfig := path.Join(home, ".gitconfig")
	if err := loadFile(globalConfig, opts); err != nil {
		return nil, err
	}

	// Load local config
	localGit, _ := git.LocalRepository().GetPath()
	localConfig := path.Join(localGit, ".git", "config")
	if err := loadFile(localConfig, opts); err != nil {
		return nil, err
	}

	// Parse args
	_, err := flags.ParseArgs(opts, os.Args)
	if err != nil {
		return nil, err
	}
	tools.Verbose = opts.Verbose
	if tools.Verbose {
		fmt.Printf("opts = %+v\n", opts)
	}

	// Return arguments
	return opts, nil
}
