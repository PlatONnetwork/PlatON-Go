// Copyright 2018-2019 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package snapshotdb

import (
	"io/ioutil"
	"os"
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

var invalidCases = []string{
	"",
	"foo",
	"foo-dx-100.log",
	".log",
	"",
	"manifest",
	"CURREN",
	"CURRENTX",
	"MANIFES",
	"MANIFEST",
	"MANIFEST-",
	"XMANIFEST-3",
	"MANIFEST-3x",
	"LOC",
	"LOCKx",
	"LO",
	"LOGx",
	"18446744073709551616.log",
	"184467440737095516150.log",
	"100",
	"100.",
	"100.lop",
}

func TestFileStorage_InvalidFileName(t *testing.T) {
	for _, name := range invalidCases {
		if _, ok := fsParseName(name); ok {
			t.Errorf("filename '%s' should be invalid", name)
		}
	}
}

func tempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "goleveldb-")
	if err != nil {
		t.Fatal(t)
	}
	t.Log("Using temp-dir:", dir)
	return dir
}

func TestFileStorage_ReadOnlyLocking(t *testing.T) {
	temp := tempDir(t)
	defer os.RemoveAll(temp)

	p1, err := openFile(temp, false)
	if err != nil {
		t.Fatal("OpenFile(1): got error: ", err)
	}

	_, err = openFile(temp, true)
	if err != nil {
		t.Logf("OpenFile(2): got error: %s (expected)", err)
	} else {
		t.Fatal("OpenFile(2): expect error")
	}

	p1.Close()

	p3, err := openFile(temp, true)
	if err != nil {
		t.Fatal("OpenFile(3): got error: ", err)
	}

	p4, err := openFile(temp, true)
	if err != nil {
		t.Fatal("OpenFile(4): got error: ", err)
	}

	_, err = openFile(temp, false)
	if err != nil {
		t.Logf("OpenFile(5): got error: %s (expected)", err)
	} else {
		t.Fatal("OpenFile(2): expect error")
	}

	p3.Close()
	p4.Close()
}
