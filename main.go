package main

import (
	"fmt"
	"os"

	flags "github.com/jessevdk/go-flags"
	"github.com/spf13/cobra"
)

type Options struct {
	Verbose bool `short:"v" long:"verbose" description:"The level of verbosity you want"`
}

var (
	VOptions = &Options{Verbose: false}
)

var parser = flags.NewParser(VOptions, flags.Default)

/*
 *
 * func loadConfigInStruct(p string, s interface{}) error {
 *   f, err := os.Open(p)
 *   if os.IsNotExist(err) {
 *     return nil
 *   }
 *   if err != nil {
 *     return err
 *   }
 *
 *   content, err := ioutil.ReadAll(f)
 *   if err != nil {
 *     return err
 *   }
 *
 *   _, err = toml.Decode(string(content), s)
 *   if err != nil {
 *     return err
 *   }
 *
 *   return nil
 * }
 *
 * func parseConfig() error {
 *   home, err := os.UserHomeDir()
 *   if err != nil {
 *     return fmt.Errorf("Unable to find home dir : %+v", err)
 *   }
 *
 *   globalConfigPath := path.Join(home, ".config", "mergo", "config")
 *   localConfigPath := path.Join(git.LocalRepository().GetPath(), ".git", "mergo.toml")
 *
 *   paths := []string{globalConfigPath, localConfigPath}
 *   structures := []interface{}{
 *     VOptions,
 *     createOptions,
 *   }
 *
 *   for _, p := range paths {
 *     for _, s := range structures {
 *       err := loadConfigInStruct(p, s)
 *       if err != nil {
 *         return fmt.Errorf("While trying to read config file : %s got error : %+v", p, err)
 *       }
 *     }
 *   }
 *   return nil
 * }
 *
 * func main() {
 *   if err := parseConfig(); err != nil {
 *     fmt.Println(err)
 *     os.Exit(1)
 *   }
 *   if _, err := parser.Parse(); err != nil {
 *     if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
 *       os.Exit(0)
 *     } else {
 *       os.Exit(1)
 *     }
 *   }
 * }
 */

func main() {
	mergo := &cobra.Command{
		Use:   "mergo",
		Short: "Mergo command",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Inside mergo command Run with args: %v\n", args)
		},
	}

	mergo.SetArgs(os.Args)
	mergo.Execute()
	fmt.Println()
}
