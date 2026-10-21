package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Print_Version(t *testing.T) {
	var version, commit, date string
	version = "dev"
	commit = "none"
	date = "unknown"

	PrintVersion()

	assert.NotNil(t, Version, Commit, Date)
	assert.Equal(t, Version, version)
	assert.Equal(t, Commit, commit)
	assert.Equal(t, Date, date)
}
