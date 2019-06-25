package main

import (
	"fmt"
	"os"
	"path"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/fatih/structs"
	"github.com/iancoleman/strcase"
	"gitlab.com/jfaucherre/mergo/git"
)

type ConfigOptions struct {
	Get        []string          `long:"get" description:"Gets a variable"`
	Set        map[string]string `long:"set" description:"Add a new variable like key:value"`
	Unset      []string          `long:"unset" description:"Unsets a variable"`
	DeleteCred []string          `long:"delete-credential" description:"Deletes the authentification credentials for the specified host"`
	Global     bool              `short:"g" long:"global" description:"Use global config file"`

	def      *structs.Struct
	confPath string
	content  map[string]interface{}
}

var configOptions = &ConfigOptions{}

// This structure contains the Options and CreateOptions merged
type fileContent struct {
	Head      string
	Base      string
	Host      string
	Remote    string
	Repo      string
	Owner     string
	Clipboard bool
	Verbose   bool
}

func init() {
	parser.AddCommand("config",
		"Manage the mergo configuration",
		`The config command helps you manage the mergo configuration
Note that the options get, set, unset and delete-credentials can be set several times to be applied several times`,
		configOptions)
}

func (me *ConfigOptions) getField(k string) (string, error) {
	// sK -> struct key
	sK := strcase.ToCamel(k)

	for _, f := range me.def.Fields() {
		if sK == f.Name() {
			switch f.Kind().String() {
			case "string":
				return me.content[k].(string), nil
			case "bool":
				return strconv.FormatBool(me.content[k].(bool)), nil
			}
		}
	}
	return "", fmt.Errorf("Field %s is not defined", k)
}

func (me *ConfigOptions) setField(k, v string) error {
	forStruct := strcase.ToCamel(k)

	for _, f := range me.def.Fields() {
		if forStruct == f.Name() {
			switch f.Kind().String() {
			case "string":
				me.content[k] = v
				break
			case "bool":
				{
					vB, err := strconv.ParseBool(v)
					if err != nil {
						return fmt.Errorf("Could not parse field %s = %s to bool\nError : %+v", k, v, err)
					}
					me.content[k] = vB
					break
				}
			}
			return nil
		}
	}
	return fmt.Errorf("Could not find field %s while trying setting it in config", k)
}

func (me *ConfigOptions) parseConf() error {
	if me.Global {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("Unable to find home dir : %+v", err)
		}
		me.confPath = path.Join(home, ".config", "mergo", "config")
	} else {
		me.confPath = path.Join(git.LocalRepository().GetPath(), ".git", "mergo.toml")
	}

	me.def = structs.New(&fileContent{})
	me.content = make(map[string]interface{})

	f, err := os.Open(me.confPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("While parsing config file %s got error %+v", me.confPath, err)
	}
	if os.IsNotExist(err) {
		return nil
	}

	_, err = toml.DecodeReader(f, &me.content)
	if err != nil {
		return fmt.Errorf("While parsing config file %s got error %+v", me.confPath, err)
	}
	return nil
}

func (me *ConfigOptions) writeConfig() error {
	f, err := os.OpenFile(me.confPath, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("While writing config got error : %+v", err)
	}

	enc := toml.NewEncoder(f)
	if err = enc.Encode(me.content); err != nil {
		return fmt.Errorf("While writing config got error : %+v", err)
	}
	return nil
}

func (me *ConfigOptions) getConfig() error {
	if len(me.Get) == 0 {
		return nil
	}

	for _, k := range me.Get {
		if _, ok := me.content[k]; !ok {
			fmt.Printf("%s = <not-set>\n", k)
		} else {
			v, _ := me.getField(k)
			fmt.Printf("%s = %+v\n", k, v)
		}
	}

	return nil
}

func (me *ConfigOptions) setConfig() error {
	for k, v := range me.Set {
		me.setField(k, v)
	}
	for _, k := range me.Unset {
		delete(me.content, k)
	}

	if err := me.writeConfig(); err != nil {
		return err
	}
	return nil
}

func (me *ConfigOptions) Execute(args []string) error {
	if err := me.parseConf(); err != nil {
		return err
	}

	if err := me.getConfig(); err != nil {
		return err
	}

	if err := me.setConfig(); err != nil {
		return err
	}

	return nil
}
