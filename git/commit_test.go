package git

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseOneCommit(t *testing.T) {
	com, err := parseOneCommit([]string{
		"commit hash (HEAD -> 166413796)",
		"Author: TotoDu38 <totodu38@gmail.com>",
		"Date:   2019-06-17 11:40:15 +0200",
		"",
		"    :sparkles: some wonderful commit message",
		"    this is some wonderful",
		"    multiline commit message",
	})
	assert.Nil(t, err)
	ti, _ := time.Parse("2006-01-02 15:04:05 -0700", "2019-06-17 11:40:15 +0200")
	assert.Equal(t, com, &Commit{
		Hash:    "hash",
		Author:  "TotoDu38 <totodu38@gmail.com>",
		Date:    ti,
		Message: ":sparkles: some wonderful commit message",
		Comment: `this is some wonderful
multiline commit message`,
	})
}

func TestParseCommitList(t *testing.T) {
	coms, err := parseCommitList(`commit hash1
Author: TotoDu38 <totodu38@gmail.com>
Date:   2019-06-17 11:40:15 +0200

    :sparkles: some wonderful commit message
    this is some wonderful
    multiline commit message

commit hash2
Author: TotoDu38 <totodu38@gmail.com>
Date:   2019-06-14 17:55:59 +0200

    another commit message
`)
	assert.Nil(t, err)
	t1, _ := time.Parse("2006-01-02 15:04:05 -0700", "2019-06-17 11:40:15 +0200")
	t2, _ := time.Parse("2006-01-02 15:04:05 -0700", "2019-06-14 17:55:59 +0200")
	assert.Equal(t, coms, []*Commit{
		&Commit{
			Hash:    "hash1",
			Author:  "TotoDu38 <totodu38@gmail.com>",
			Date:    t1,
			Message: ":sparkles: some wonderful commit message",
			Comment: `this is some wonderful
multiline commit message`,
		},
		&Commit{
			Hash:    "hash2",
			Author:  "TotoDu38 <totodu38@gmail.com>",
			Date:    t2,
			Message: "another commit message",
			Comment: "",
		},
	})
}
