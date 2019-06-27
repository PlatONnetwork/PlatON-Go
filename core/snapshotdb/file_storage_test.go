package snapshotdb

import (
	"testing"
)

func Test_fsParseName(t *testing.T) {
	names := []string{
		"0000000020-0x09bde6c3ba69637b526b9e05909617ceb9821881d67bdcd263bf591e45e4f846.log",
		"current",
	}
	for _, name := range names {
		fd, ok := fsParseName(name)
		if !ok {
			t.Error("must right name:", name)
		}
		if fd.String() != name {
			t.Error("must the same name:", name)
		}
	}
	wrongNames := []string{
		"cccccc",
	}
	for _, name := range wrongNames {
		_, ok := fsParseName(name)
		if ok {
			t.Error("must wrong name:", name)
		}
	}
}
