package operations

import (
	"testing"

	"github.com/OffchainLabs/prysm/v6/consensus-types/blocks"
	"github.com/OffchainLabs/prysm/v6/consensus-types/interfaces"
	ethpb "github.com/OffchainLabs/prysm/v6/proto/prysm/v1alpha1"
	"github.com/OffchainLabs/prysm/v6/runtime/version"
	common "github.com/OffchainLabs/prysm/v6/testing/spectest/shared/common/operations"
	"github.com/OffchainLabs/prysm/v6/testing/util"
)

func blockWithAttesterSlashing(asSSZ []byte) (interfaces.SignedBeaconBlock, error) {
	as := &ethpb.AttesterSlashing{}
	if err := as.UnmarshalSSZ(asSSZ); err != nil {
		return nil, err
	}
	b := util.NewBeaconBlockCapella()
	b.Block.Body = &ethpb.BeaconBlockBodyCapella{AttesterSlashings: []*ethpb.AttesterSlashing{as}}
	return blocks.NewSignedBeaconBlock(b)
}

func RunAttesterSlashingTest(t *testing.T, config string) {
	common.RunAttesterSlashingTest(t, config, version.String(version.Capella), blockWithAttesterSlashing, sszToState)
}
