package semver_test

import (
	"testing"

	"github.com/Nidal-Bakir/go-semver"
	"github.com/stretchr/testify/assert"
)

func TestParseSemVerFunction(t *testing.T) {
	a := assert.New(t)

	ver, err := semver.ParseSemVer("1.0.0")
	a.NoError(err)
	a.Equal("1.0.0", ver.String())
	a.Equal(1, ver.Major)
	a.Equal(0, ver.Minor)
	a.Equal(0, ver.Patch)
	a.Equal("", ver.PreRelease)
	a.Equal("", ver.BuildMetadata)

}
