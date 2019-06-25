package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/BurntSushi/toml"
	flags "github.com/jessevdk/go-flags"
	"gitlab.com/jfaucherre/mergo/git"
)

type Options struct {
	Verbose bool `short:"v" long:"verbose" description:"The level of verbosity you want"`
}

var (
	VOptions = &Options{Verbose: false}
)

var parser = flags.NewParser(VOptions, flags.Default)

func addConfig(p string) error {
	f, err := os.Open(p)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	_, err = toml.Decode(string(content), VOptions)
	if err != nil {
		return err
	}

	_, err = toml.Decode(string(content), VCreateOptions)
	if err != nil {
		return err
	}

	return nil
}

func parseConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("Unable to find home dir : %+v", err)
	}

	globalConfigPath := path.Join(home, ".config", "mergo", "config")

	localConfigPath := path.Join(git.LocalRepository().GetPath(), ".git", "mergo.toml")

	for _, p := range []string{globalConfigPath, localConfigPath} {
		err := addConfig(p)
		if err != nil {
			return fmt.Errorf("While trying to read config file : %s got error : %+v", p, err)
		}
	}
	return nil
}

func main() {
	if err := parseConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if _, err := parser.Parse(); err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}
}
