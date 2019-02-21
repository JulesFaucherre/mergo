package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoteRegexp(t *testing.T) {
	testValues := []struct {
		Value    string
		Expected [][]string
	}{
		{
			Value: "https://gitlab.com/jfaucherre/mergo.git",
			Expected: [][]string{
				{
					"https://gitlab.com/jfaucherre/mergo.git",
					"gitlab.com",
					"jfaucherre",
					"mergo",
				},
				nil,
			},
		},
		{
			Value: "git@gitlab.com:jfaucherre/mergo.git",
			Expected: [][]string{
				nil,
				{
					"git@gitlab.com:jfaucherre/mergo.git",
					"gitlab.com",
					"jfaucherre",
					"mergo",
				},
			},
		},
	}

	for _, testValue := range testValues {
		assert.Equal(t, testValue.Expected[0], httpsR.FindStringSubmatch(testValue.Value))
		assert.Equal(t, testValue.Expected[1], sshR.FindStringSubmatch(testValue.Value))
	}
}
