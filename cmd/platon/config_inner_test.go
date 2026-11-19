package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInnerAdminKeyFileJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "platon.json")
	if err := os.WriteFile(path, []byte(`{"Eth":{"Miner":{"InnerAdminKeyFile":"/secure/platon/inner-admin.key"}}}`), 0600); err != nil {
		t.Fatal(err)
	}
	var cfg platonConfig
	if err := loadConfigFile(path, &cfg); err != nil {
		t.Fatal(err)
	}
	if got := cfg.Eth.Miner.InnerAdminKeyFile; got != "/secure/platon/inner-admin.key" {
		t.Fatalf("Miner.InnerAdminKeyFile = %q", got)
	}
}
