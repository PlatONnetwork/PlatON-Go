package tx

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"math/big"
	"testing"
)

func TestTx(t *testing.T) {
	tx_json_str := `{
        "tx_type": 1,
        "input":[
            {
                "ephemeral_pk": "566a57f8dac3ce75177ad9f47cca5c8f7c159a19b7dffa7199a08edc43b8d777",
                "sign_pk": "be73cefa7f638a03d2e14ad7ec04fbb1906a84579cfb8176c2faa599534ac36f",
                "quantity": 1,
                "blinding": "0100000000000000000000000000000000000000000000000000000000000000",
                "view_sk": "0100000000000000000000000000000000000000000000000000000000000000",
                "spend_sk": "0100000000000000000000000000000000000000000000000000000000000000"
            },
            {
                "ephemeral_pk": "fc6d651a50268da1726e5ecb75c0a8926cc91694f29df9d2cd22ce8ebcf54213",
                "sign_pk": "0affa2249e8200e43709001363b9dce05289e0f1436630ba7bb04924438e9862",
                "quantity": 4,
                "blinding": "0100000000000000000000000000000000000000000000000000000000000000",
                "view_sk": "0400000000000000000000000000000000000000000000000000000000000000",
                "spend_sk": "0400000000000000000000000000000000000000000000000000000000000000"
            }
        ],
        "output": [
            {
                "quantity": 4,
                "view_pk": "0affa2249e8200e43709001363b9dce05289e0f1436630ba7bb04924438e9862",
                "spend_pk": "0affa2249e8200e43709001363b9dce05289e0f1436630ba7bb04924438e9862"
            },
            {
                "quantity": 1,
                "view_pk": "0affa2249e8200e43709001363b9dce05289e0f1436630ba7bb04924438e9862",
                "spend_pk": "0affa2249e8200e43709001363b9dce05289e0f1436630ba7bb04924438e9862"
            }
        ],
        "authorized_address": "0affa2249e8200e43709001363b9dce05289e0f1"
    }`
	proof, err := CreateConfidentialTx([]byte(tx_json_str))
	if err != nil {
		t.Error(err)
	}
	t.Log(len(proof))

	result, err := VerifyConfidentialTx(proof)
	if err != nil {
		t.Error(err)
	}

	t.Log(len(result))
	t.Log(hexutil.Encode(result))
	//var tx TxLog
	//if err = rlp.DecodeBytes(result, &tx); err != nil {
	//	t.Error(err)
	//}
	tx := TxLog{TxType: 1, Inputs: []Note{
		Note{EphemeralPk: common.BigToHash(big.NewInt(1)).Bytes(),
			SpendingPk: common.BigToHash(big.NewInt(1)).Bytes(),
			Token:      common.BigToHash(big.NewInt(1)).Bytes()},
	},
		Outputs: []Note{Note{EphemeralPk: common.BigToHash(big.NewInt(1)).Bytes(),
			SpendingPk: common.BigToHash(big.NewInt(1)).Bytes(),
			Token:      common.BigToHash(big.NewInt(1)).Bytes()},
		}}
	s, _ := rlp.EncodeToBytes(&tx)
	t.Log(hexutil.Encode(s))
}
