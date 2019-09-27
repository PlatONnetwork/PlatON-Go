package types

import (
	"errors"

	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

// EncodeExtra encode cbft version and `QuorumCert` as extra data.
func EncodeExtra(cbftVersion byte, qc *QuorumCert) ([]byte, error) {
	extra := []byte{cbftVersion}
	bxBytes, err := rlp.EncodeToBytes(qc)
	if err != nil {
		return nil, err
	}
	extra = append(extra, bxBytes...)
	return extra, nil
}

// DecodeExtra decode extra data as cbft version and `QuorumCert`.
func DecodeExtra(extra []byte) (byte, *QuorumCert, error) {
	if len(extra) == 0 {
		return 0, nil, errors.New("empty extra")
	}
	version := extra[0]
	var qc QuorumCert
	err := rlp.DecodeBytes(extra[1:], &qc)
	if err != nil {
		return 0, nil, err
	}
	return version, &qc, nil
}
