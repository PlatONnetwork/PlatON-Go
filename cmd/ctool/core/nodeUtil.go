// Copyright 2021 The PlatON Network Authors
// This file is part of PlatON-Go.
//
// PlatON-Go is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// PlatON-Go is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with PlatON-Go. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	//todo
	from = `0xc1f330b214668beac2e6418dd651b09c759a4bf5`
)

var (
	dir, _      = os.Getwd()
	abiFilePath = "../test/contracta.cpp.abi.json"
	//configPath  = "../config.json"
	configPath = filepath.Join(dir, "../config.json")
	pkFilePath = "../test/privateKeys.txt"
)

func parseConfig(t *testing.T) {
	err := parseConfigJson(configPath)
	assert.Nil(t, err, fmt.Sprintf("%v", err))
}
