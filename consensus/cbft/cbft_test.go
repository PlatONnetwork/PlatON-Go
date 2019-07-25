package cbft

import (
	"crypto/ecdsa"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/validator"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestThreshold(t *testing.T) {
	f := &Cbft{}
	assert.Equal(t, 1, f.threshold(1))
	assert.Equal(t, 2, f.threshold(2))
	assert.Equal(t, 3, f.threshold(3))
	assert.Equal(t, 3, f.threshold(4))
	assert.Equal(t, 4, f.threshold(5))
	assert.Equal(t, 5, f.threshold(6))
	assert.Equal(t, 5, f.threshold(7))

}

func generateKeys(num int) ([]*ecdsa.PrivateKey, []*bls.SecretKey) {
	pk := make([]*ecdsa.PrivateKey, 0)
	sk := make([]*bls.SecretKey, 0)

	for i := 0; i < num; i++ {
		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		ecdsaKey, _ := crypto.GenerateKey()
		pk = append(pk, ecdsaKey)
		sk = append(sk, &blsKey)
	}
	return pk, sk
}
func TestBls(t *testing.T) {
	bls.Init(bls.CurveFp254BNb)
	num := 4
	pk, sk := generateKeys(num)
	owner := sk[0]
	nodes := make([]params.CbftNode, num)
	for i := 0; i < num; i++ {

		nodes[i].Node = *discover.NewNode(discover.PubkeyID(&pk[i].PublicKey), nil, 0, 0)
		nodes[i].BlsPubKey = *sk[i].GetPublicKey()
	}

	agency := validator.NewStaticAgency(nodes)

	cbft := &Cbft{
		validatorPool: validator.NewValidatorPool(agency, 0, nodes[0].Node.ID),
		config: types.Config{
			Option: &types.OptionsConfig{
				BlsPriKey: owner,
			},
		},
	}

	pb := &protocols.PrepareVote{}
	cbft.signMsgByBls(pb)
	msg, _ := pb.CannibalizeBytes()
	assert.True(t, cbft.validatorPool.Verify(0, 0, msg, pb.Sign()))
}
