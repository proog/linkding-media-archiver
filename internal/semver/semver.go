package semver

import (
	"cmp"
	"fmt"
	"strconv"
	"strings"
)

func Parse(version string) (Semver, error) {
	parts := strings.SplitN(version, ".", 3)

	if len(parts) != 3 {
		return Semver{}, fmt.Errorf("Failed to parse \"%s\" as a version string", version)
	}

	major, majorErr := strconv.Atoi(parts[0])
	minor, minorErr := strconv.Atoi(parts[1])

	patchStr, suffix, _ := strings.Cut(parts[2], "-")
	patch, patchErr := strconv.Atoi(patchStr)

	if majorErr != nil || minorErr != nil || patchErr != nil {
		return Semver{}, fmt.Errorf("Failed to parse \"%s\" as a version string", version)
	}

	semver := Semver{major, minor, patch, suffix}
	return semver, nil
}

func Compare(x, y Semver) int {
	majorCmp := cmp.Compare(x.Major, y.Major)
	if majorCmp != 0 {
		return majorCmp
	}

	minorCmp := cmp.Compare(x.Minor, y.Minor)
	if minorCmp != 0 {
		return minorCmp
	}

	patchCmp := cmp.Compare(x.Patch, y.Patch)
	if patchCmp != 0 {
		return patchCmp
	}

	return 0
}

func (v Semver) String() string {
	str := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)

	if v.Suffix != "" {
		str = str + fmt.Sprintf("-%s", v.Suffix)
	}

	return str
}
