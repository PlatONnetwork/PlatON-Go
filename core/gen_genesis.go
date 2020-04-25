package core

import (
	"encoding/json"
	"errors"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/common/math"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

var _ = (*genesisSpecMarshaling)(nil)

// MarshalJSON marshals as JSON.
func (g Genesis) MarshalJSON() ([]byte, error) {
	type Genesis struct {
		Config        *params.ChainConfig               `json:"config"`
		EconomicModel *xcom.EconomicModel               `json:"economicModel"`
		Nonce         hexutil.Bytes                     `json:"nonce"`
		Timestamp     math.HexOrDecimal64               `json:"timestamp"`
		ExtraData     hexutil.Bytes                     `json:"extraData"`
		GasLimit      math.HexOrDecimal64               `json:"gasLimit"   gencodec:"required"`
		Coinbase      common.Address                    `json:"coinbase"`
		Alloc         map[common.Address]GenesisAccount `json:"alloc"      gencodec:"required"`
		Number        math.HexOrDecimal64               `json:"number"`
		GasUsed       math.HexOrDecimal64               `json:"gasUsed"`
		ParentHash    common.Hash                       `json:"parentHash"`
	}
	var enc Genesis
	enc.Config = g.Config
	enc.EconomicModel = g.EconomicModel
	enc.Nonce = g.Nonce
	enc.Timestamp = math.HexOrDecimal64(g.Timestamp)
	enc.ExtraData = g.ExtraData
	enc.GasLimit = math.HexOrDecimal64(g.GasLimit)
	enc.Coinbase = g.Coinbase
	enc.Alloc = g.Alloc
	enc.Number = math.HexOrDecimal64(g.Number)
	enc.GasUsed = math.HexOrDecimal64(g.GasUsed)
	enc.ParentHash = g.ParentHash
	return json.Marshal(&enc)
}

// UnmarshalJSON unmarshals from JSON.
func (g *Genesis) UnmarshalJSON(input []byte) error {
	type Genesis struct {
		Config        *params.ChainConfig               `json:"config"`
		EconomicModel *xcom.EconomicModel               `json:"economicModel"`
		Nonce         *hexutil.Bytes                    `json:"nonce"`
		Timestamp     *math.HexOrDecimal64              `json:"timestamp"`
		ExtraData     *hexutil.Bytes                    `json:"extraData"`
		GasLimit      *math.HexOrDecimal64              `json:"gasLimit"   gencodec:"required"`
		Coinbase      *common.Address                   `json:"coinbase"`
		Alloc         map[common.Address]GenesisAccount `json:"alloc"      gencodec:"required"`
		Number        *math.HexOrDecimal64              `json:"number"`
		GasUsed       *math.HexOrDecimal64              `json:"gasUsed"`
		ParentHash    *common.Hash                      `json:"parentHash"`
	}
	var dec Genesis
	dec.EconomicModel = g.EconomicModel
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Config != nil {
		g.Config = dec.Config
	}
	if dec.EconomicModel != nil {
		g.EconomicModel = dec.EconomicModel
	}
	if dec.Nonce != nil {
		g.Nonce = *dec.Nonce
	}
	if dec.Timestamp != nil {
		g.Timestamp = uint64(*dec.Timestamp)
	}
	if dec.ExtraData != nil {
		g.ExtraData = *dec.ExtraData
	}
	if dec.GasLimit == nil {
		return errors.New("missing required field 'gasLimit' for Genesis")
	}
	g.GasLimit = uint64(*dec.GasLimit)
	if dec.Coinbase != nil {
		g.Coinbase = *dec.Coinbase
	}
	if dec.Alloc == nil {
		return errors.New("missing required field 'alloc' for Genesis")
	} else {
		g.Alloc = dec.Alloc
	}
	if dec.Number != nil {
		g.Number = uint64(*dec.Number)
	}
	if dec.GasUsed != nil {
		g.GasUsed = uint64(*dec.GasUsed)
	}
	if dec.ParentHash != nil {
		g.ParentHash = *dec.ParentHash
	}
	return nil
}
