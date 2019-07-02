package create

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gitlab.com/jfaucherre/mergo/models"
)

var (
	CreateCmd = &cobra.Command{
		Use:   "create",
		Short: "The command to create a pull request",
		Long: `This command will get the informations you give it from the parameters
    and extract the others from your git configuration to create a
    pull request`,
		Run: run,
	}
	options = &models.CreateOptions{}
)

var (
	httpsR = regexp.MustCompile(`https://(\w+\.\w+)/([\w-]+)/([\w-]+).git`)
	sshR   = regexp.MustCompile(`git@(\w+\.\w+):([\w-]+)/([\w-]+).git`)
)

func init() {
	cobra.OnInitialize(initConfig)
	CreateCmd.Flags().StringVarP(&options.Head, "head", "d", "", "The head branch you want to merge into the base")
	CreateCmd.Flags().StringVarP(&options.Base, "base", "b", "master", "The base branch you want to merge into")
	CreateCmd.Flags().StringVar(&options.Host, "host", "", "The git host you use, ie github, gitlab, etc.")

	CreateCmd.Flags().StringVar(&options.Remote, "remote", "origin", "The remote to use")
	CreateCmd.Flags().StringVar(&options.Repo, "repository", "", "The name of the repository on which you want to make the pull request")
	CreateCmd.Flags().StringVar(&options.Owner, "owner", "", "The owner of the repository")
	CreateCmd.Flags().BoolVarP(&options.Clipboard, "copy-clipboard", "c", false, "Copies the merge request adress to the clipboard")

	viper.BindPFlag("base", CreateCmd.PersistentFlags().Lookup("base"))
	viper.BindPFlag("clipboard", CreateCmd.PersistentFlags().Lookup("copy-clipboard"))
}

func loadConfig(dir, fname string) {
	fullPath := path.Join(dir, fname)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		ioutil.WriteFile(fullPath, []byte{}, 0644)
	}

	viper.SetConfigType("toml")
	viper.SetConfigFile(fullPath)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Can't read config:", err)
		os.Exit(1)
	}
}

func initConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	globalDir := path.Join(home, ".config", "mergo")
	loadConfig(globalDir, "mergorc.toml")
}
