package tools

import (
	"errors"
	"reflect"
	"strings"

	url "github.com/whilp/git-urls"
	"gitlab.com/jfaucherre/mergo/credentials"
	"gitlab.com/jfaucherre/mergo/logger"
)

var (
	ErrInvalidRepositoryURL = errors.New("invalid given url")
)

type Host interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Name() string
}

func LoadHostCredentials(host Host) error {
	creds, err := credentials.GetHostConfig(host.Name())
	if err != nil {
		return err
	}

	return host.Unmarshal(creds)
}

func SaveHostCredentials(host Host) error {
	creds, err := host.Marshal()
	if err != nil {
		return err
	}

	return credentials.WriteHostConfig(host.Name(), creds)
}

// Takes a host and get its credentials from the user. To define how the
// credentials must be asked, use the struct tags parse-name, parse-info and
// parse-reg:
//  - parse-name is the name of value that will asked, typically the user will
//    be asked a question like "please enter your <parse-name>"
//    parse-name is required
//  - parse-info is just a sentence adding more info to your request
//  - parse-reg is the regexp that the user input must match to be valid
func AskForHostCredentials(host Host) error {
	if host == nil {
		return ErrNoNil
	}

	// if we got a ptr or an interface we get the value behind
	rv := reflect.ValueOf(host)
	for rv.Kind() == reflect.Ptr || rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}

	if !rv.IsValid() {
		return ErrInvalidType
	}
	if rv.Kind() != reflect.Struct {
		return ErrExpectedStruct
	}

	t := rv.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fT := field.Type
		if fT.Kind() != reflect.Slice || fT.Elem().Kind() != reflect.Uint8 {
			return ErrUnhandledFieldType
		}

		name := field.Tag.Get("parse-name")
		info := field.Tag.Get("parse-info")
		reg := field.Tag.Get("parse-reg")
		logger.Info("parsing field %s with name %s, info %s and regexp %s\n", field.Name, name, info, reg)
		value, err := queryUserInfo(name, info, reg)
		if err != nil {
			return err
		}
		rv.Field(i).SetBytes(value)
	}

	keep, err := AskYesNo("Do you want these credentials to be kept for your next merge requests ?")
	if err != nil {
		return err
	}

	if !keep {
		return nil
	}

	content, err := host.Marshal()
	if err != nil {
		return nil
	}
	if err := credentials.WriteHostConfig(host.Name(), content); err != nil {
		return err
	}
	return nil
}

func RepoInfoFromURL(urlString string) (result struct {
	Owner string
	Repo  string
}, err error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return
	}
	path := u.Path

	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	if strings.HasSuffix(path, ".git") {
		path = path[:len(path)-4]
	}

	splitted := strings.Split(path, "/")
	if len(splitted) != 2 {
		err = ErrInvalidRepositoryURL
		return
	}
	result.Owner = splitted[0]
	result.Repo = splitted[1]
	return
}
