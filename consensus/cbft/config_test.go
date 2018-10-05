package cbft

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	SetConfigFile("D:/workspace/golang/src/Platon-go/config/cbft-config.toml")
	fmt.Println(Config().sealers[1].PublicKey)
}
