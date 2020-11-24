package git

import (
	"io"
	"io/ioutil"
	"os"
	"path"

	ini "gopkg.in/ini.v1"
)

func (me *Repo) ConfigPath() string {
	return path.Join(me.path, ".git", "config")
}

// LoadConfig takes any object and try to unmarshal the config into it
// use the ini tag to make it parsable
func (me *Repo) LoadConfig(obj interface{}) error {
	var files []string

	xdg := os.Getenv("XDG_CONFIG_HOME")
	if xdg != "" {
		files = append(files, path.Join(xdg, "git", "config"))
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	files = append(files,
		path.Join(home, ".gitconfig"),
		path.Join(home, ".config", "git", "config"),
	)

	files = append(files, me.ConfigPath())

	for _, f := range files {
		r, err := os.Open(f)
		if os.IsNotExist(err) {
			continue
		}
		if err := readConfigFile(obj, r); err != nil {
			return err
		}
	}
	return nil
}

func readConfigFile(obj interface{}, r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	cfg, err := ini.Load(data)
	if err != nil {
		return err
	}

	err = cfg.MapTo(obj)
	if err == os.ErrNotExist {
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

// GetEditor returns the editor used by git to edit text
func (me *Repo) GetEditor() (string, error) {
	obj := new(struct {
		Core struct {
			Editor string `ini:"editor"`
		} `ini:"core"`
	})
	err := me.LoadConfig(obj)
	if err != nil {
		return "", err
	}

	editor := obj.Core.Editor
	if editor != "" {
		return editor, nil
	}
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor, nil
	}
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor, nil
	}

	return "vi", nil
}
