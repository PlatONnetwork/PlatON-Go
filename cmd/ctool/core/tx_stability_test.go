package core

import (
	"testing"
)

func TestPrepareAccount(t *testing.T) {
	//parseConfigJson(configPath)
	//err := PrepareAccount(10, pkFilePath, "0xDE0B6B3A7640000")
	//if err != nil {
	//	t.Fatalf(err.Error())
	//}
}

func TestStressTest(t *testing.T) {
	parseConfigJson(configPath)
	err := StabilityTest(pkFilePath, 1, 10)
	if err != nil {
		t.Fatalf(err.Error())
	}
}
