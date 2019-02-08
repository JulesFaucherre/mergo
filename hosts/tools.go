package hosts

import "strings"

func ownerAndRepoSSH(remote string) (string, string) {
	splitted := strings.FieldsFunc(remote, func(r rune) bool {
		return r == ':' || r == '/' || r == '.'
	})

	return splitted[2], splitted[3]
}

func ownerAndRepoHTTP(remote string) (string, string) {
	splitted := strings.FieldsFunc(remote, func(r rune) bool {
		return r == '/' || r == '.'
	})

	return splitted[3], splitted[4]
}

func ownerAndRepo(remote string) (string, string) {
	if strings.HasPrefix(remote, "http") {
		return ownerAndRepoHTTP(remote)
	}
	return ownerAndRepoSSH(remote)
}
