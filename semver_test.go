package semver_test

import (
	"fmt"
	"math/rand"
	"slices"
	"testing"

	"github.com/Nidal-Bakir/go-semver"
	"github.com/stretchr/testify/assert"
)

func TestParseSemVerFunction(t *testing.T) {
	a := assert.New(t)

	type testCase struct {
		semverStr      string
		expectedSemver semver.SemVer
	}

	var testData = []testCase{
		testCase{
			semverStr:      "1.0.0-0.3.7",
			expectedSemver: semver.SemVer{Major: 1, Minor: 0, Patch: 0, PreRelease: "0.3.7"},
		},
		testCase{
			semverStr:      "1.0.0-alpha",
			expectedSemver: semver.SemVer{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha"},
		},
		testCase{
			semverStr:      "1.0.0-alpha+001",
			expectedSemver: semver.SemVer{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha", BuildMetadata: "001"},
		},
		testCase{
			semverStr:      "1.0.0-alpha.1",
			expectedSemver: semver.SemVer{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha.1"},
		},

		testCase{
			semverStr:      "1.0.0-alpha.beta",
			expectedSemver: semver.SemVer{Major: 1, Minor: 0, Patch: 0, PreRelease: "alpha.beta"},
		},
		testCase{
			semverStr:      "1.0.0-beta",
			expectedSemver: semver.SemVer{Major: 1, Minor: 0, Patch: 0, PreRelease: "beta"},
		},
		testCase{
			semverStr:      "1.0.0-beta+exp.sha.5114f85",
			expectedSemver: semver.SemVer{Major: 1, Minor: 0, Patch: 0, PreRelease: "beta", BuildMetadata: "exp.sha.5114f85"},
		},
		testCase{
			semverStr:      "1.0.0-beta.2",
			expectedSemver: semver.SemVer{Major: 1, Minor: 0, Patch: 0, PreRelease: "beta.2"},
		},
		testCase{
			semverStr:      "1.0.0-beta.11",
			expectedSemver: semver.SemVer{Major: 1, Minor: 0, Patch: 0, PreRelease: "beta.11"},
		},
		testCase{
			semverStr:      "1.0.0-rc.1",
			expectedSemver: semver.SemVer{Major: 1, Minor: 0, Patch: 0, PreRelease: "rc.1"},
		},
		testCase{
			semverStr:      "1.0.0-x.7.z.92",
			expectedSemver: semver.SemVer{Major: 1, Minor: 0, Patch: 0, PreRelease: "x.7.z.92"},
		},
		testCase{
			semverStr:      "1.0.0+21AF26D3----117B344092BD",
			expectedSemver: semver.SemVer{Major: 1, Minor: 0, Patch: 0, BuildMetadata: "21AF26D3----117B344092BD"},
		},
		testCase{
			semverStr:      "1.0.0+20130313144700",
			expectedSemver: semver.SemVer{Major: 1, Minor: 0, Patch: 0, BuildMetadata: "20130313144700"},
		},
		testCase{
			semverStr:      "1.0.0",
			expectedSemver: semver.SemVer{Major: 1, Minor: 0, Patch: 0},
		},
		testCase{
			semverStr:      "2.0.0",
			expectedSemver: semver.SemVer{Major: 2, Minor: 0, Patch: 0},
		},
		testCase{
			semverStr:      "11.11.11",
			expectedSemver: semver.SemVer{Major: 11, Minor: 11, Patch: 11},
		},
		testCase{
			semverStr:      "62.99.57962",
			expectedSemver: semver.SemVer{Major: 62, Minor: 99, Patch: 57962},
		},
		testCase{
			semverStr:      "17897798789984631.14461777.1100000",
			expectedSemver: semver.SemVer{Major: 17897798789984631, Minor: 14461777, Patch: 1100000},
		},
	}

	for i := range testData {
		semverStr := testData[i].semverStr
		expectedSemver := testData[i].expectedSemver
		ver, err := semver.Parse(semverStr)
		a.NoError(err)
		a.True(semver.IsValid(semverStr))
		a.Equal(expectedSemver.String(), ver.String())
		a.Equal(expectedSemver.Major, ver.Major)
		a.Equal(expectedSemver.Minor, ver.Minor)
		a.Equal(expectedSemver.Patch, ver.Patch)
		a.Equal(expectedSemver.PreRelease, ver.PreRelease)
		a.Equal(expectedSemver.BuildMetadata, ver.BuildMetadata)
		a.Equal(expectedSemver.IsPreRelease(), ver.IsPreRelease())
		a.True(expectedSemver.IsEquql(ver))

		if i+1 != len(testData) {

			nextSemverStr := testData[i+1].semverStr
			nextSemver, err := semver.Parse(nextSemverStr)
			a.NoError(err)
			a.True(ver.IsLessOrEquql(nextSemver), fmt.Sprint("(", ver.String(), ") should be less then of equal (", nextSemver.String(), ")"))
			a.False(ver.IsGrater(nextSemver), fmt.Sprint("(", ver.String(), ") should be less then of equal (", nextSemver.String(), ")"))

			a.True(semver.IsValid(nextSemverStr))

			cmpResult, err := semver.Compare(semverStr, nextSemverStr)
			a.NoError(err)
			a.True(cmpResult <= 0)
		}
	}
}

func TestSemverPrecedenceWithBuildMetadata(t *testing.T) {
	a := assert.New(t)

	semver1, err := semver.Parse("1.0.0")
	a.NoError(err)
	semver2, err := semver.Parse("1.0.0+1231456")
	a.NoError(err)
	a.True(semver1.IsEquql(semver2))
}

func TestSemverPrecedenceWithPreRelease(t *testing.T) {
	a := assert.New(t)

	type testCase struct {
		semver1 string
		semver2 string
	}

	var testData = []testCase{
		testCase{
			semver1: "1.0.0-alpha",
			semver2: "1.0.0-alpha.1",
		},
		testCase{
			semver1: "1.0.0-alpha",
			semver2: "1.0.0-alpha.beta",
		},
		testCase{
			semver1: "1.0.0-rc.1",
			semver2: "1.0.0",
		},
	}

	for i := range testData {
		semver1, err := semver.Parse(testData[i].semver1)
		a.NoError(err)
		semver2, err := semver.Parse(testData[i].semver2)
		a.NoError(err)

		a.False(semver1.IsEquql(semver2))
		a.True(semver1.IsLess(semver2))
		a.True(semver1.IsLessOrEquql(semver2))
		a.True(semver2.IsGrater(semver1))
		a.True(semver2.IsGraterOrEquql(semver1))
	}

}

func TestParseErrorsAndInValid(t *testing.T) {
	a := assert.New(t)

	testData := []string{
		"1..0",
		"1..",
		"1.0.",
		"1.0",
		"1.0-b",
		"1.0-b",
		"1b.0.0",
	}

	for _, v := range testData {
		_, err := semver.Parse(v)
		a.Error(err)
		a.False(semver.IsValid(v))
	}
}

func TestSortAndSortStr(t *testing.T) {
	a := assert.New(t)

	expectedTestData := []string{
		"1.0.0-0.3.7",
		"1.0.0-alpha",
		"1.0.0-alpha.1",
		"1.0.0-alpha.beta",
		"1.0.0-beta",
		"1.0.0-beta.2",
		"1.0.0-beta.11",
		"1.0.0-rc.1",
		"1.0.0-x.7.z.92",
		"1.0.0",
		"2.0.0",
		"11.11.11",
		"62.99.57962",
	}

	randomizedTestData := slices.Clone(expectedTestData)
	randomizedTestDataSemver := make([]semver.SemVer, len(randomizedTestData))

	for i := range randomizedTestData {
		j := rand.Intn(i + 1)
		randomizedTestData[i], randomizedTestData[j] = randomizedTestData[j], randomizedTestData[i]
	}

	for i := range randomizedTestData {
		randomizedTestDataSemver[i] = semver.MustParse(randomizedTestData[i])
	}

	err := semver.SortStr(randomizedTestData)
	a.NoError(err)

	semver.Sort(randomizedTestDataSemver)
	a.NoError(err)

	for i, expected := range expectedTestData {
		randSemverObj := randomizedTestDataSemver[i].String()
		randSemverStr := randomizedTestData[i]
		a.True(randSemverObj == randSemverStr && randSemverStr == expected, fmt.Sprint(randSemverObj, " & ", randSemverStr, " & ", expected, " should be equal"))
	}
}
