package semver

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	semver, err := Parse("1.44.2-rc.4")

	if err != nil {
		t.Fatal(err)
	}

	if semver.Major != 1 {
		t.Errorf("Expectd major version 1, got %d", semver.Major)
	}
	if semver.Minor != 44 {
		t.Errorf("Expectd minor version 44, got %d", semver.Minor)
	}
	if semver.Patch != 2 {
		t.Errorf("Expectd patch version 2, got %d", semver.Patch)
	}
	if semver.Suffix != "rc.4" {
		t.Errorf("Expectd suffix rc.4, got %s", semver.Suffix)
	}
}

func TestCompare(t *testing.T) {
	tests := []struct {
		x        Semver
		y        Semver
		expected int
	}{
		{Semver{Major: 1, Minor: 43, Patch: 5}, Semver{Major: 1, Minor: 44, Patch: 0}, -1},
		{Semver{Major: 1, Minor: 43, Patch: 5}, Semver{Major: 1, Minor: 43, Patch: 5, Suffix: "ignored"}, 0},
		{Semver{Major: 1, Minor: 44}, Semver{Major: 1, Minor: 44}, 0},
		{Semver{Major: 1, Minor: 43, Patch: 12}, Semver{Major: 1, Minor: 43, Patch: 9}, 1},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s | %s", test.x, test.y), func(t *testing.T) {
			result := Compare(test.x, test.y)

			if result != test.expected {
				t.Errorf("Expected comparison result to be %d, got %d", test.expected, result)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		version  Semver
		expected string
	}{
		{Semver{Major: 1, Minor: 44, Patch: 5}, "1.44.5"},
		{Semver{Major: 1, Minor: 44, Patch: 5, Suffix: "alpha.4"}, "1.44.5-alpha.4"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			actual := fmt.Sprintf("%s", test.version)

			if actual != test.expected {
				t.Errorf("Expected version to be formatted as %s, got %s", test.expected, actual)
			}
		})
	}
}
