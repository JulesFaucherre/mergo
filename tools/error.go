package tools

import (
	"gitlab.com/jfaucherre/mergo/logger"
)

func CheckError(err error) {
	if err != nil {
		logger.Fatal("%+v", err)
	}
}
