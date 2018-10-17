package cbft

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"path/filepath"
	"sync"
)

var configFile string

func SetConfigFile(confFile string) {
	configFile = confFile
}

type cbftConfig struct {
	Period            uint64 `json:"period"` // Number of seconds between blocks to enforce
	Epoch             uint64 `json:"epoch"`  // Epoch length to reset votes and checkpoint
	MaxNetworkLatency uint32 `json:"maxNetworkLatency"`
	//Sealers []common.Address `json:"sealers"`
	Sealers []string `json:"sealers"`
}

var (
	cfg     *cbftConfig
	once    sync.Once
	cfgLock = new(sync.RWMutex)
)

func Config() *cbftConfig {
	once.Do(ReloadConfig)
	cfgLock.RLock()
	defer cfgLock.RUnlock()
	return cfg
}

func ReloadConfig() {
	filePath, err := filepath.Abs(configFile)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Printf("parse toml file once. filePath: %s\n", filePath)
	config := new(cbftConfig)
	if _, err := toml.DecodeFile(filePath, config); err != nil {
		fmt.Println(err)
		panic(err)
	}
	cfgLock.Lock()
	defer cfgLock.Unlock()
	cfg = config
}
