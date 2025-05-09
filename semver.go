// Semantic Versioning 2.0.0
//
// see https://semver.org
package semver

import (
	"cmp"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

var (
	ErrInvalidSemVerSyntax = errors.New("invalid semver syntax")
)

type SemVer struct {
	Major         int
	Minor         int
	Patch         int
	PreRelease    string
	BuildMetadata string
}

func (s SemVer) IsPreRelease() bool {
	return len(s.PreRelease) != 0
}

func (s SemVer) String() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprint(s.Major, ".", s.Minor, ".", s.Patch))

	if len(s.PreRelease) != 0 {
		b.WriteString("-")
		b.WriteString(s.PreRelease)
	}

	if len(s.BuildMetadata) != 0 {
		b.WriteString("+")
		b.WriteString(s.BuildMetadata)
	}

	return b.String()
}

func (s SemVer) IsGraterOrEquql(o SemVer) bool {
	cmpResults := s.generateCompareToOtherResult(o)
	for _, result := range cmpResults {
		if result == 0 {
			continue
		}
		return result > 0
	}
	return true
}

func (s SemVer) IsGrater(o SemVer) bool {
	cmpResults := s.generateCompareToOtherResult(o)
	for _, result := range cmpResults {
		if result == 0 {
			continue
		}
		return result > 0
	}
	return false
}

func (s SemVer) IsLess(o SemVer) bool {
	cmpResults := s.generateCompareToOtherResult(o)
	for _, result := range cmpResults {
		if result == 0 {
			continue
		}
		return result < 0
	}
	return false
}

func (s SemVer) IsLessOrEquql(o SemVer) bool {
	cmpResults := s.generateCompareToOtherResult(o)
	for _, result := range cmpResults {
		if result == 0 {
			continue
		}
		return result < 0
	}
	return true
}

func (s SemVer) generateCompareToOtherResult(o SemVer) []int {
	res := make([]int, 4)
	res[0] = cmp.Compare(s.Major, o.Major)
	res[1] = cmp.Compare(s.Minor, o.Minor)
	res[2] = cmp.Compare(s.Patch, o.Patch)
	res[3] = s.comparePreRelease(o)
	return res
}

func (s SemVer) IsEquql(o SemVer) bool {
	return s.Major == o.Major && s.Minor == o.Minor && s.Patch == o.Patch && s.PreRelease == o.PreRelease
}

// Compare returns
//
//	-1 if this is less than other,
//	 0 if this equals other,
//	+1 if this is greater than other.
func (s SemVer) Compare(other SemVer) int {
	cmpResults := s.generateCompareToOtherResult(other)
	for _, result := range cmpResults {
		if result == 0 {
			continue
		}
		return result
	}
	return 0
}

// When major, minor, and patch are equal, a pre-release version has lower precedence than a normal version:
//
// Example: 1.0.0-alpha < 1.0.0.
//
// Precedence for two pre-release versions with the same major, minor,
// and patch version MUST be determined by comparing each dot separated
// identifier from left to right until a difference is found as follows:
//
//  1. Identifiers consisting of only digits are compared numerically.
//
//  2. Identifiers with letters or hyphens are compared lexically in ASCII sort order.
//
//  3. Numeric identifiers always have lower precedence than non-numeric identifiers.
//
//  4. A larger set of pre-release fields has a higher precedence than a smaller set, if all of the preceding identifiers are equal.
//
// Example: 1.0.0-alpha < 1.0.0-alpha.1 < 1.0.0-alpha.beta < 1.0.0-beta < 1.0.0-beta.2 < 1.0.0-beta.11 < 1.0.0-rc.1 < 1.0.0.
func (s SemVer) comparePreRelease(o SemVer) int {
	thisLen := len(s.PreRelease)
	otherLen := len(o.PreRelease)
	if thisLen == 0 || otherLen == 0 {
		if thisLen == 0 && otherLen == 0 {
			return 0
		}
		if thisLen == 0 {
			return 1
		}
		return -1
	}

	thisPreReleaseSplit := strings.Split(s.PreRelease, ".")
	otherPreReleaseSplit := strings.Split(o.PreRelease, ".")
	for i := range min(len(thisPreReleaseSplit), len(otherPreReleaseSplit)) {
		thisPart := thisPreReleaseSplit[i]
		otherPart := otherPreReleaseSplit[i]

		tDigit, tOk := mayParseDigit(thisPart)
		oDigit, oOk := mayParseDigit(otherPart)

		if tOk || oOk { // one part is a digit
			if tOk && oOk { // the two parts are digits
				if tDigit == oDigit {
					continue
				}
				return cmp.Compare(tDigit, oDigit)
			}

			// one digit and the other is not
			if tOk {
				return -1
			}
			return 1
		}

		// the two parts are strings
		cmpResult := cmp.Compare(thisPart, otherPart)
		if cmpResult == 0 {
			continue
		}
		return cmpResult
	}

	// The smallest split completed at this point is equal to the other split.
	// Return the comparison of the lengths of the two PreRelease versions.
	return cmp.Compare(thisLen, otherLen)
}

func mayParseDigit(s string) (int, bool) {
	v, err := strconv.Atoi(s)
	return v, err == nil
}

func MustParse(semverStr string) SemVer {
	v, err := Parse(semverStr)
	if err != nil {
		panic(err)
	}
	return v
}

func Parse(semverStr string) (SemVer, error) {
	parts := make([]strings.Builder, 5)
	partIndex := 0
	didEnterPreReleasePart := false
	didEnterBuildMetadataPart := false
	for _, c := range semverStr {
		if c == '.' && partIndex < 2 {
			partIndex++
			continue
		}

		if c == '-' && !didEnterPreReleasePart && !didEnterBuildMetadataPart {
			didEnterPreReleasePart = true
			partIndex = 3
			continue
		}

		if c == '+' && !didEnterBuildMetadataPart {
			didEnterBuildMetadataPart = true
			partIndex = 4
			continue
		}

		parts[partIndex].WriteRune(c)
	}

	var semver SemVer

	major, err := strconv.Atoi(parts[0].String())
	if err != nil {
		return semver, ErrInvalidSemVerSyntax
	}
	semver.Major = major

	minor, err := strconv.Atoi(parts[1].String())
	if err != nil {
		return semver, ErrInvalidSemVerSyntax
	}
	semver.Minor = minor

	patch, err := strconv.Atoi(parts[2].String())
	if err != nil {
		return semver, ErrInvalidSemVerSyntax
	}
	semver.Patch = patch

	semver.PreRelease = parts[3].String()
	semver.BuildMetadata = parts[4].String()

	return semver, nil
}

func IsValid(v string) bool {
	_, err := Parse(v)
	return err == nil
}

func Compare(v1, v2 string) (int, error) {
	parsedV1, err := Parse(v1)
	if err != nil {
		return 0, err
	}
	parsedV2, err := Parse(v2)
	if err != nil {
		return 0, err
	}
	return parsedV1.Compare(parsedV2), nil
}

// Sort sorts a list of semantic versions in a ascending order
func Sort(slice []SemVer) {
	sort.Sort(semverSlice(slice))
}

// semverSlice implements [sort.Interface] for sorting semantic versions.
type semverSlice []SemVer

func (sv semverSlice) Len() int           { return len(sv) }
func (sv semverSlice) Swap(i, j int)      { sv[i], sv[j] = sv[j], sv[i] }
func (sv semverSlice) Less(i, j int) bool { return sv[i].IsLessOrEquql(sv[j]) }

// Sort sorts a list of string semantic versions in a ascending order.
func SortStr(slice []string) error {
	m := make(semverSlice, len(slice))
	
	for i := range slice {
		v, err := Parse(slice[i])
		if err != nil {
			return err
		}
		m[i] = v
	}
	
	sort.Sort(m)
	
	for i := range slice {
		slice[i] = m[i].String()
	}
	
	return nil
}
