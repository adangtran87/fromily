package main

import (
	"fmt"
	"strconv"
)

type Version struct {
	Major uint64
	Minor uint64
	Patch uint64
}

func (v *Version) String() string {
	major_str := strconv.FormatUint(v.Major, 10)
	minor_str := strconv.FormatUint(v.Minor, 10)
	patch_str := strconv.FormatUint(v.Patch, 10)
	return fmt.Sprintf("%s.%s.%s", major_str, minor_str, patch_str)
}

var Ver = Version{
	0,
	0,
	1,
}
