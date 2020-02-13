package xcom

import "testing"

func TestGetDefaultEMConfig(t *testing.T) {
	if getDefaultEMConfig(DefaultMainNet) == nil {
		t.Error("DefaultMainNet can't be nil config")
	}
	if getDefaultEMConfig(DefaultTestNet) == nil {
		t.Error("DefaultTestNet can't be nil config")
	}
	if getDefaultEMConfig(DefaultDemoNet) == nil {
		t.Error("DefaultDemoNet can't be nil config")
	}
	if getDefaultEMConfig(10) == nil {
		t.Error("DefaultDemoNet can't be nil config")
	}
}
