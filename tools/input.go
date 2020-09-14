package tools

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"regexp"

	"gitlab.com/jfaucherre/mergo/logger"
)

var (
	ErrNoNil              = errors.New("unexpected nil value")
	ErrInvalidType        = errors.New("invalid type")
	ErrExpectedStruct     = errors.New("expected struct type")
	ErrNotHandledType     = errors.New("type not handled")
	ErrUnhandledFieldType = errors.New("only handles slice of byte as field type")
)

type GetLineParams struct {
	Question []byte
	Regexp   *regexp.Regexp
}

func GetLine(params GetLineParams) (ret []byte, err error) {
	reader := bufio.NewReader(os.Stdin)
	logger.Debug("getting user line with params\n\tquestion: %s\n\treg: %+v\n", string(params.Question), params.Regexp)

	for {
		fmt.Println(string(params.Question))
		ret, err = reader.ReadBytes('\n')
		if err != nil {
			return
		}
		if params.Regexp != nil && params.Regexp.Match(ret) {
			fmt.Printf("Invalid input: %s, must match: %s\n", ret, params.Regexp)
			continue
		}

		break
	}
	// remove \n
	ret = ret[:len(ret)-1]

	return
}

func AskYesNo(s string) (bool, error) {
	res, err := GetLine(GetLineParams{
		Question: []byte(s + " ([y]/n)"),
		Regexp:   regexp.MustCompile("^[yn]?$"),
	})
	if err != nil {
		return false, err
	}
	res = bytes.TrimSpace(res)
	notAgreed := len(res) == 1 && res[0] == 'n'
	return !notAgreed, nil
}

func queryUserInfo(name, info, reg string) ([]byte, error) {
	question := fmt.Sprintf("Please input your %s", name)
	if len(info) > 0 {
		question = question + "\n" + info
	}
	var r *regexp.Regexp
	if len(reg) > 0 {
		r = regexp.MustCompile(reg)
	}

	return GetLine(GetLineParams{
		Question: []byte(question),
		Regexp:   r,
	})
}
