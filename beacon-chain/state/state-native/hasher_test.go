package state_native_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/prysmaticlabs/go-bitfield"
	statenative "github.com/prysmaticlabs/prysm/v5/beacon-chain/state/state-native"
	"github.com/prysmaticlabs/prysm/v5/config/params"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/blocks"
	"github.com/prysmaticlabs/prysm/v5/encoding/bytesutil"
	enginev1 "github.com/prysmaticlabs/prysm/v5/proto/engine/v1"
	ethpb "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v5/testing/assert"
	"github.com/prysmaticlabs/prysm/v5/testing/require"
	"github.com/prysmaticlabs/prysm/v5/testing/util"
)

func TestComputeFieldRootsWithHasher_Phase0(t *testing.T) {
	beaconState, err := util.NewBeaconState(util.FillRootsNaturalOpt)
	require.NoError(t, err)
	require.NoError(t, beaconState.SetGenesisTime(123))
	require.NoError(t, beaconState.SetGenesisValidatorsRoot(genesisValidatorsRoot()))
	require.NoError(t, beaconState.SetSlot(123))
	require.NoError(t, beaconState.SetFork(fork()))
	require.NoError(t, beaconState.SetLatestBlockHeader(latestBlockHeader()))
	historicalRoots, err := util.PrepareRoots(int(params.BeaconConfig().SlotsPerHistoricalRoot))
	require.NoError(t, err)
	require.NoError(t, beaconState.SetHistoricalRoots(historicalRoots))
	require.NoError(t, beaconState.SetEth1Data(eth1Data()))
	require.NoError(t, beaconState.SetEth1DataVotes([]*ethpb.Eth1Data{eth1Data()}))
	require.NoError(t, beaconState.SetEth1DepositIndex(123))
	require.NoError(t, beaconState.SetValidators([]*ethpb.Validator{validator()}))
	require.NoError(t, beaconState.SetBalances([]uint64{1, 2, 3}))
	randaoMixes, err := util.PrepareRoots(int(params.BeaconConfig().EpochsPerHistoricalVector))
	require.NoError(t, err)
	require.NoError(t, beaconState.SetRandaoMixes(randaoMixes))
	require.NoError(t, beaconState.SetSlashings([]uint64{1, 2, 3}))
	require.NoError(t, beaconState.AppendPreviousEpochAttestations(pendingAttestation("previous")))
	require.NoError(t, beaconState.AppendCurrentEpochAttestations(pendingAttestation("current")))
	require.NoError(t, beaconState.SetJustificationBits(justificationBits()))
	require.NoError(t, beaconState.SetPreviousJustifiedCheckpoint(checkpoint("previous")))
	require.NoError(t, beaconState.SetCurrentJustifiedCheckpoint(checkpoint("current")))
	require.NoError(t, beaconState.SetFinalizedCheckpoint(checkpoint("finalized")))

	nativeState, ok := beaconState.(*statenative.BeaconState)
	require.Equal(t, true, ok)
	protoState, ok := nativeState.ToProtoUnsafe().(*ethpb.BeaconState)
	require.Equal(t, true, ok)

	initState, err := statenative.InitializeFromProtoPhase0(protoState)
	require.NoError(t, err)
	s, ok := initState.(*statenative.BeaconState)
	require.Equal(t, true, ok)
	root, err := statenative.ComputeFieldRootsWithHasher(context.Background(), s)
	require.NoError(t, err)
	expected := [][]byte{
		{0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0x67, 0x76, 0x72, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0x58, 0xba, 0xf, 0x9b, 0x4f, 0x63, 0x1c, 0xa6, 0x19, 0xb1, 0xa2, 0x1f, 0xd1, 0x29, 0xc7, 0x67, 0x9c, 0x32, 0x4, 0x1f, 0xcf, 0x4e, 0x64, 0x9b, 0x8f, 0x21, 0xb4, 0xe6, 0xa5, 0xc9, 0xc, 0x38},
		{0x8b, 0x5, 0x59, 0x78, 0xed, 0xbe, 0x2c, 0xde, 0xa6, 0xf, 0x52, 0xdc, 0x16, 0x83, 0xa0, 0x5d, 0x8, 0xc3, 0x37, 0x91, 0x3a, 0xf6, 0xfa, 0x6, 0x62, 0xc9, 0x6, 0xb1, 0x41, 0x48, 0xaf, 0xec},
		{0xcc, 0xa1, 0x69, 0xa8, 0x8b, 0x1e, 0x49, 0x5, 0x1f, 0xe4, 0xd8, 0xff, 0x82, 0xa2, 0x2e, 0xf0, 0x54, 0xd1, 0x13, 0xc9, 0x8e, 0xb9, 0x82, 0xa6, 0x9e, 0x42, 0x2, 0xec, 0x97, 0x6f, 0x33, 0x88},
		{0xcc, 0xa1, 0x69, 0xa8, 0x8b, 0x1e, 0x49, 0x5, 0x1f, 0xe4, 0xd8, 0xff, 0x82, 0xa2, 0x2e, 0xf0, 0x54, 0xd1, 0x13, 0xc9, 0x8e, 0xb9, 0x82, 0xa6, 0x9e, 0x42, 0x2, 0xec, 0x97, 0x6f, 0x33, 0x88},
		{0xf, 0xad, 0xd3, 0x92, 0x4b, 0xda, 0xfa, 0xc6, 0x61, 0x50, 0xb7, 0xdf, 0x5f, 0x2c, 0xd0, 0x94, 0xc3, 0xaf, 0x41, 0x9d, 0xa, 0xea, 0x50, 0x96, 0x82, 0x62, 0x1c, 0x72, 0x26, 0x20, 0x6b, 0xac},
		{0xc9, 0x4e, 0x2c, 0xb0, 0x20, 0xe3, 0xe7, 0x8c, 0x5c, 0xbd, 0xeb, 0x9b, 0xa5, 0x7b, 0x53, 0x50, 0xca, 0xfe, 0xe9, 0x48, 0x9e, 0x8d, 0xf8, 0x4a, 0xe6, 0x8d, 0x9c, 0x97, 0x81, 0x74, 0xb, 0x5e},
		{0x71, 0x52, 0xd2, 0x9b, 0x87, 0x3c, 0x8a, 0xd9, 0x51, 0x55, 0xc0, 0x42, 0xb, 0xc4, 0x12, 0xa4, 0x79, 0xf5, 0x7d, 0x37, 0x16, 0xf4, 0x90, 0x72, 0x5d, 0xe0, 0x34, 0xb4, 0x2, 0x8c, 0x39, 0xe4},
		{0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0xf, 0xf7, 0x4f, 0xe1, 0xa9, 0x72, 0x9c, 0x95, 0xf0, 0xe1, 0xde, 0xa4, 0x32, 0xc, 0x67, 0x52, 0x23, 0x13, 0x9e, 0xe2, 0x40, 0x8d, 0xf6, 0x18, 0x57, 0xf0, 0x1a, 0x4a, 0xad, 0x46, 0xce, 0x42},
		{0xf9, 0x11, 0x2c, 0xc2, 0x71, 0x70, 0xde, 0x47, 0x26, 0xeb, 0x26, 0xd4, 0xa4, 0xe8, 0x68, 0xb, 0x16, 0xa2, 0x6e, 0x52, 0x54, 0xe, 0x5c, 0x83, 0x17, 0x3, 0xea, 0xdd, 0xd5, 0xa7, 0xb2, 0x3f},
		{0x64, 0xbd, 0x40, 0xa7, 0x10, 0x44, 0x84, 0xed, 0xf3, 0x5f, 0xc3, 0x5d, 0x7b, 0xbe, 0xe8, 0x75, 0xbf, 0x66, 0xcb, 0xce, 0x77, 0xfa, 0x0, 0x3, 0xdd, 0xfb, 0x80, 0xd2, 0x77, 0x1b, 0xc2, 0x8},
		{0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0x54, 0xd3, 0xce, 0x8a, 0x3f, 0xfd, 0x21, 0x3a, 0xb4, 0xa6, 0xd, 0xb, 0x9f, 0xf2, 0x88, 0xf0, 0xb1, 0x44, 0x9d, 0xb1, 0x2, 0x95, 0x67, 0xdf, 0x6f, 0x28, 0xa9, 0x68, 0xcd, 0xaa, 0x8c, 0x54},
		{0xeb, 0x8, 0xb4, 0x1b, 0x76, 0xa2, 0x23, 0xbb, 0x4a, 0xd3, 0x78, 0xca, 0x2e, 0xe8, 0x2c, 0xa1, 0xbf, 0x45, 0xf2, 0x58, 0xdf, 0x39, 0xdf, 0x43, 0x40, 0xb, 0x96, 0xcf, 0xfd, 0x9a, 0x87, 0x85},
		{0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0x76, 0x88, 0xa0, 0x68, 0x45, 0x25, 0x8f, 0xd5, 0xf9, 0xb2, 0xb0, 0x42, 0x68, 0x6b, 0x51, 0xcc, 0x29, 0x94, 0x63, 0x85, 0xec, 0xf5, 0x47, 0xf0, 0x9c, 0x46, 0x86, 0xa9, 0x99, 0x7d, 0x29, 0x6c},
		{0x41, 0x44, 0x52, 0xff, 0x8c, 0xa6, 0xb3, 0x2e, 0xcc, 0x5e, 0x63, 0x8f, 0x8e, 0x7d, 0xe7, 0x52, 0x42, 0x94, 0x55, 0x2f, 0x89, 0xdd, 0x1e, 0x3c, 0xb0, 0xf4, 0x51, 0x51, 0x36, 0x81, 0x72, 0x1},
		{0xa9, 0xbb, 0x6a, 0x1f, 0x5d, 0x86, 0x7d, 0xa7, 0x5a, 0x7d, 0x9d, 0x8d, 0xc0, 0x15, 0xb7, 0x0, 0xee, 0xa9, 0x68, 0x51, 0x88, 0x57, 0x5a, 0xd9, 0x4e, 0x1d, 0x8e, 0x44, 0xbf, 0xdc, 0x73, 0xff},
	}
	assert.DeepEqual(t, expected, root)
}

func TestComputeFieldRootsWithHasher_Altair(t *testing.T) {
	beaconState, err := util.NewBeaconStateAltair(util.FillRootsNaturalOptAltair)
	require.NoError(t, err)
	require.NoError(t, beaconState.SetGenesisTime(123))
	require.NoError(t, beaconState.SetGenesisValidatorsRoot(genesisValidatorsRoot()))
	require.NoError(t, beaconState.SetSlot(123))
	require.NoError(t, beaconState.SetFork(fork()))
	require.NoError(t, beaconState.SetLatestBlockHeader(latestBlockHeader()))
	historicalRoots, err := util.PrepareRoots(int(params.BeaconConfig().SlotsPerHistoricalRoot))
	require.NoError(t, err)
	require.NoError(t, beaconState.SetHistoricalRoots(historicalRoots))
	require.NoError(t, beaconState.SetEth1Data(eth1Data()))
	require.NoError(t, beaconState.SetEth1DataVotes([]*ethpb.Eth1Data{eth1Data()}))
	require.NoError(t, beaconState.SetEth1DepositIndex(123))
	require.NoError(t, beaconState.SetValidators([]*ethpb.Validator{validator()}))
	require.NoError(t, beaconState.SetBalances([]uint64{1, 2, 3}))
	randaoMixes, err := util.PrepareRoots(int(params.BeaconConfig().EpochsPerHistoricalVector))
	require.NoError(t, err)
	require.NoError(t, beaconState.SetRandaoMixes(randaoMixes))
	require.NoError(t, beaconState.SetSlashings([]uint64{1, 2, 3}))
	require.NoError(t, beaconState.SetPreviousParticipationBits([]byte{1, 2, 3}))
	require.NoError(t, beaconState.SetCurrentParticipationBits([]byte{1, 2, 3}))
	require.NoError(t, beaconState.SetJustificationBits(justificationBits()))
	require.NoError(t, beaconState.SetPreviousJustifiedCheckpoint(checkpoint("previous")))
	require.NoError(t, beaconState.SetCurrentJustifiedCheckpoint(checkpoint("current")))
	require.NoError(t, beaconState.SetFinalizedCheckpoint(checkpoint("finalized")))
	require.NoError(t, beaconState.SetInactivityScores([]uint64{1, 2, 3}))
	require.NoError(t, beaconState.SetCurrentSyncCommittee(syncCommittee("current")))
	require.NoError(t, beaconState.SetNextSyncCommittee(syncCommittee("next")))

	nativeState, ok := beaconState.(*statenative.BeaconState)
	require.Equal(t, true, ok)
	protoState, ok := nativeState.ToProtoUnsafe().(*ethpb.BeaconStateAltair)
	require.Equal(t, true, ok)
	initState, err := statenative.InitializeFromProtoAltair(protoState)
	require.NoError(t, err)
	s, ok := initState.(*statenative.BeaconState)
	require.Equal(t, true, ok)

	root, err := statenative.ComputeFieldRootsWithHasher(context.Background(), s)
	require.NoError(t, err)
	expected := [][]byte{
		{0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0x67, 0x76, 0x72, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0x58, 0xba, 0xf, 0x9b, 0x4f, 0x63, 0x1c, 0xa6, 0x19, 0xb1, 0xa2, 0x1f, 0xd1, 0x29, 0xc7, 0x67, 0x9c, 0x32, 0x4, 0x1f, 0xcf, 0x4e, 0x64, 0x9b, 0x8f, 0x21, 0xb4, 0xe6, 0xa5, 0xc9, 0xc, 0x38},
		{0x8b, 0x5, 0x59, 0x78, 0xed, 0xbe, 0x2c, 0xde, 0xa6, 0xf, 0x52, 0xdc, 0x16, 0x83, 0xa0, 0x5d, 0x8, 0xc3, 0x37, 0x91, 0x3a, 0xf6, 0xfa, 0x6, 0x62, 0xc9, 0x6, 0xb1, 0x41, 0x48, 0xaf, 0xec},
		{0xcc, 0xa1, 0x69, 0xa8, 0x8b, 0x1e, 0x49, 0x5, 0x1f, 0xe4, 0xd8, 0xff, 0x82, 0xa2, 0x2e, 0xf0, 0x54, 0xd1, 0x13, 0xc9, 0x8e, 0xb9, 0x82, 0xa6, 0x9e, 0x42, 0x2, 0xec, 0x97, 0x6f, 0x33, 0x88},
		{0xcc, 0xa1, 0x69, 0xa8, 0x8b, 0x1e, 0x49, 0x5, 0x1f, 0xe4, 0xd8, 0xff, 0x82, 0xa2, 0x2e, 0xf0, 0x54, 0xd1, 0x13, 0xc9, 0x8e, 0xb9, 0x82, 0xa6, 0x9e, 0x42, 0x2, 0xec, 0x97, 0x6f, 0x33, 0x88},
		{0xf, 0xad, 0xd3, 0x92, 0x4b, 0xda, 0xfa, 0xc6, 0x61, 0x50, 0xb7, 0xdf, 0x5f, 0x2c, 0xd0, 0x94, 0xc3, 0xaf, 0x41, 0x9d, 0xa, 0xea, 0x50, 0x96, 0x82, 0x62, 0x1c, 0x72, 0x26, 0x20, 0x6b, 0xac},
		{0xc9, 0x4e, 0x2c, 0xb0, 0x20, 0xe3, 0xe7, 0x8c, 0x5c, 0xbd, 0xeb, 0x9b, 0xa5, 0x7b, 0x53, 0x50, 0xca, 0xfe, 0xe9, 0x48, 0x9e, 0x8d, 0xf8, 0x4a, 0xe6, 0x8d, 0x9c, 0x97, 0x81, 0x74, 0xb, 0x5e},
		{0x71, 0x52, 0xd2, 0x9b, 0x87, 0x3c, 0x8a, 0xd9, 0x51, 0x55, 0xc0, 0x42, 0xb, 0xc4, 0x12, 0xa4, 0x79, 0xf5, 0x7d, 0x37, 0x16, 0xf4, 0x90, 0x72, 0x5d, 0xe0, 0x34, 0xb4, 0x2, 0x8c, 0x39, 0xe4},
		{0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0xf, 0xf7, 0x4f, 0xe1, 0xa9, 0x72, 0x9c, 0x95, 0xf0, 0xe1, 0xde, 0xa4, 0x32, 0xc, 0x67, 0x52, 0x23, 0x13, 0x9e, 0xe2, 0x40, 0x8d, 0xf6, 0x18, 0x57, 0xf0, 0x1a, 0x4a, 0xad, 0x46, 0xce, 0x42},
		{0xf9, 0x11, 0x2c, 0xc2, 0x71, 0x70, 0xde, 0x47, 0x26, 0xeb, 0x26, 0xd4, 0xa4, 0xe8, 0x68, 0xb, 0x16, 0xa2, 0x6e, 0x52, 0x54, 0xe, 0x5c, 0x83, 0x17, 0x3, 0xea, 0xdd, 0xd5, 0xa7, 0xb2, 0x3f},
		{0x64, 0xbd, 0x40, 0xa7, 0x10, 0x44, 0x84, 0xed, 0xf3, 0x5f, 0xc3, 0x5d, 0x7b, 0xbe, 0xe8, 0x75, 0xbf, 0x66, 0xcb, 0xce, 0x77, 0xfa, 0x0, 0x3, 0xdd, 0xfb, 0x80, 0xd2, 0x77, 0x1b, 0xc2, 0x8},
		{0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0x75, 0xb2, 0xae, 0x1d, 0xd8, 0xca, 0xe6, 0x4d, 0xa8, 0xc5, 0xc9, 0x19, 0x8, 0x96, 0xaf, 0x9b, 0xe6, 0xf6, 0x99, 0xb9, 0x58, 0x56, 0x5b, 0x25, 0xea, 0x9c, 0x86, 0x5e, 0x96, 0x6a, 0x48, 0xb},
		{0x75, 0xb2, 0xae, 0x1d, 0xd8, 0xca, 0xe6, 0x4d, 0xa8, 0xc5, 0xc9, 0x19, 0x8, 0x96, 0xaf, 0x9b, 0xe6, 0xf6, 0x99, 0xb9, 0x58, 0x56, 0x5b, 0x25, 0xea, 0x9c, 0x86, 0x5e, 0x96, 0x6a, 0x48, 0xb},
		{0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0x76, 0x88, 0xa0, 0x68, 0x45, 0x25, 0x8f, 0xd5, 0xf9, 0xb2, 0xb0, 0x42, 0x68, 0x6b, 0x51, 0xcc, 0x29, 0x94, 0x63, 0x85, 0xec, 0xf5, 0x47, 0xf0, 0x9c, 0x46, 0x86, 0xa9, 0x99, 0x7d, 0x29, 0x6c},
		{0x41, 0x44, 0x52, 0xff, 0x8c, 0xa6, 0xb3, 0x2e, 0xcc, 0x5e, 0x63, 0x8f, 0x8e, 0x7d, 0xe7, 0x52, 0x42, 0x94, 0x55, 0x2f, 0x89, 0xdd, 0x1e, 0x3c, 0xb0, 0xf4, 0x51, 0x51, 0x36, 0x81, 0x72, 0x1},
		{0xa9, 0xbb, 0x6a, 0x1f, 0x5d, 0x86, 0x7d, 0xa7, 0x5a, 0x7d, 0x9d, 0x8d, 0xc0, 0x15, 0xb7, 0x0, 0xee, 0xa9, 0x68, 0x51, 0x88, 0x57, 0x5a, 0xd9, 0x4e, 0x1d, 0x8e, 0x44, 0xbf, 0xdc, 0x73, 0xff},
		{0xf9, 0x11, 0x2c, 0xc2, 0x71, 0x70, 0xde, 0x47, 0x26, 0xeb, 0x26, 0xd4, 0xa4, 0xe8, 0x68, 0xb, 0x16, 0xa2, 0x6e, 0x52, 0x54, 0xe, 0x5c, 0x83, 0x17, 0x3, 0xea, 0xdd, 0xd5, 0xa7, 0xb2, 0x3f},
		{0x3d, 0xf3, 0x66, 0xd4, 0x12, 0x40, 0x3f, 0x28, 0xeb, 0xe4, 0x19, 0x59, 0xae, 0xab, 0x4d, 0xf3, 0x98, 0x88, 0x7f, 0x1e, 0x58, 0xa, 0x5d, 0xd4, 0xeb, 0xe5, 0x5d, 0x3d, 0x11, 0x70, 0x24, 0x76},
		{0xd6, 0x4c, 0xb1, 0xac, 0x61, 0x7, 0x26, 0xbb, 0xd3, 0x27, 0x2a, 0xcd, 0xdd, 0x55, 0xf, 0x2b, 0x6a, 0xe8, 0x1, 0x31, 0x48, 0x66, 0x2f, 0x98, 0x7b, 0x6d, 0x27, 0x69, 0xd9, 0x40, 0xcc, 0x37},
	}
	assert.DeepEqual(t, expected, root)
}

func TestComputeFieldRootsWithHasher_Bellatrix(t *testing.T) {
	beaconState, err := util.NewBeaconStateBellatrix(util.FillRootsNaturalOptBellatrix)
	require.NoError(t, err)
	require.NoError(t, beaconState.SetGenesisTime(123))
	require.NoError(t, beaconState.SetGenesisValidatorsRoot(genesisValidatorsRoot()))
	require.NoError(t, beaconState.SetSlot(123))
	require.NoError(t, beaconState.SetFork(fork()))
	require.NoError(t, beaconState.SetLatestBlockHeader(latestBlockHeader()))
	historicalRoots, err := util.PrepareRoots(int(params.BeaconConfig().SlotsPerHistoricalRoot))
	require.NoError(t, err)
	require.NoError(t, beaconState.SetHistoricalRoots(historicalRoots))
	require.NoError(t, beaconState.SetEth1Data(eth1Data()))
	require.NoError(t, beaconState.SetEth1DataVotes([]*ethpb.Eth1Data{eth1Data()}))
	require.NoError(t, beaconState.SetEth1DepositIndex(123))
	require.NoError(t, beaconState.SetValidators([]*ethpb.Validator{validator()}))
	require.NoError(t, beaconState.SetBalances([]uint64{1, 2, 3}))
	randaoMixes, err := util.PrepareRoots(int(params.BeaconConfig().EpochsPerHistoricalVector))
	require.NoError(t, err)
	require.NoError(t, beaconState.SetRandaoMixes(randaoMixes))
	require.NoError(t, beaconState.SetSlashings([]uint64{1, 2, 3}))
	require.NoError(t, beaconState.SetPreviousParticipationBits([]byte{1, 2, 3}))
	require.NoError(t, beaconState.SetCurrentParticipationBits([]byte{1, 2, 3}))
	require.NoError(t, beaconState.SetJustificationBits(justificationBits()))
	require.NoError(t, beaconState.SetPreviousJustifiedCheckpoint(checkpoint("previous")))
	require.NoError(t, beaconState.SetCurrentJustifiedCheckpoint(checkpoint("current")))
	require.NoError(t, beaconState.SetFinalizedCheckpoint(checkpoint("finalized")))
	require.NoError(t, beaconState.SetInactivityScores([]uint64{1, 2, 3}))
	require.NoError(t, beaconState.SetCurrentSyncCommittee(syncCommittee("current")))
	require.NoError(t, beaconState.SetNextSyncCommittee(syncCommittee("next")))
	wrappedHeader, err := blocks.WrappedExecutionPayloadHeader(executionPayloadHeader())
	require.NoError(t, err)
	require.NoError(t, beaconState.SetLatestExecutionPayloadHeader(wrappedHeader))

	nativeState, ok := beaconState.(*statenative.BeaconState)
	require.Equal(t, true, ok)
	protoState, ok := nativeState.ToProtoUnsafe().(*ethpb.BeaconStateBellatrix)
	require.Equal(t, true, ok)
	initState, err := statenative.InitializeFromProtoBellatrix(protoState)
	require.NoError(t, err)
	s, ok := initState.(*statenative.BeaconState)
	require.Equal(t, true, ok)

	root, err := statenative.ComputeFieldRootsWithHasher(context.Background(), s)
	require.NoError(t, err)
	expected := [][]byte{
		{0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0x67, 0x76, 0x72, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0x58, 0xba, 0xf, 0x9b, 0x4f, 0x63, 0x1c, 0xa6, 0x19, 0xb1, 0xa2, 0x1f, 0xd1, 0x29, 0xc7, 0x67, 0x9c, 0x32, 0x4, 0x1f, 0xcf, 0x4e, 0x64, 0x9b, 0x8f, 0x21, 0xb4, 0xe6, 0xa5, 0xc9, 0xc, 0x38},
		{0x8b, 0x5, 0x59, 0x78, 0xed, 0xbe, 0x2c, 0xde, 0xa6, 0xf, 0x52, 0xdc, 0x16, 0x83, 0xa0, 0x5d, 0x8, 0xc3, 0x37, 0x91, 0x3a, 0xf6, 0xfa, 0x6, 0x62, 0xc9, 0x6, 0xb1, 0x41, 0x48, 0xaf, 0xec},
		{0xcc, 0xa1, 0x69, 0xa8, 0x8b, 0x1e, 0x49, 0x5, 0x1f, 0xe4, 0xd8, 0xff, 0x82, 0xa2, 0x2e, 0xf0, 0x54, 0xd1, 0x13, 0xc9, 0x8e, 0xb9, 0x82, 0xa6, 0x9e, 0x42, 0x2, 0xec, 0x97, 0x6f, 0x33, 0x88},
		{0xcc, 0xa1, 0x69, 0xa8, 0x8b, 0x1e, 0x49, 0x5, 0x1f, 0xe4, 0xd8, 0xff, 0x82, 0xa2, 0x2e, 0xf0, 0x54, 0xd1, 0x13, 0xc9, 0x8e, 0xb9, 0x82, 0xa6, 0x9e, 0x42, 0x2, 0xec, 0x97, 0x6f, 0x33, 0x88},
		{0xf, 0xad, 0xd3, 0x92, 0x4b, 0xda, 0xfa, 0xc6, 0x61, 0x50, 0xb7, 0xdf, 0x5f, 0x2c, 0xd0, 0x94, 0xc3, 0xaf, 0x41, 0x9d, 0xa, 0xea, 0x50, 0x96, 0x82, 0x62, 0x1c, 0x72, 0x26, 0x20, 0x6b, 0xac},
		{0xc9, 0x4e, 0x2c, 0xb0, 0x20, 0xe3, 0xe7, 0x8c, 0x5c, 0xbd, 0xeb, 0x9b, 0xa5, 0x7b, 0x53, 0x50, 0xca, 0xfe, 0xe9, 0x48, 0x9e, 0x8d, 0xf8, 0x4a, 0xe6, 0x8d, 0x9c, 0x97, 0x81, 0x74, 0xb, 0x5e},
		{0x71, 0x52, 0xd2, 0x9b, 0x87, 0x3c, 0x8a, 0xd9, 0x51, 0x55, 0xc0, 0x42, 0xb, 0xc4, 0x12, 0xa4, 0x79, 0xf5, 0x7d, 0x37, 0x16, 0xf4, 0x90, 0x72, 0x5d, 0xe0, 0x34, 0xb4, 0x2, 0x8c, 0x39, 0xe4},
		{0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0xf, 0xf7, 0x4f, 0xe1, 0xa9, 0x72, 0x9c, 0x95, 0xf0, 0xe1, 0xde, 0xa4, 0x32, 0xc, 0x67, 0x52, 0x23, 0x13, 0x9e, 0xe2, 0x40, 0x8d, 0xf6, 0x18, 0x57, 0xf0, 0x1a, 0x4a, 0xad, 0x46, 0xce, 0x42},
		{0xf9, 0x11, 0x2c, 0xc2, 0x71, 0x70, 0xde, 0x47, 0x26, 0xeb, 0x26, 0xd4, 0xa4, 0xe8, 0x68, 0xb, 0x16, 0xa2, 0x6e, 0x52, 0x54, 0xe, 0x5c, 0x83, 0x17, 0x3, 0xea, 0xdd, 0xd5, 0xa7, 0xb2, 0x3f},
		{0x64, 0xbd, 0x40, 0xa7, 0x10, 0x44, 0x84, 0xed, 0xf3, 0x5f, 0xc3, 0x5d, 0x7b, 0xbe, 0xe8, 0x75, 0xbf, 0x66, 0xcb, 0xce, 0x77, 0xfa, 0x0, 0x3, 0xdd, 0xfb, 0x80, 0xd2, 0x77, 0x1b, 0xc2, 0x8},
		{0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0x75, 0xb2, 0xae, 0x1d, 0xd8, 0xca, 0xe6, 0x4d, 0xa8, 0xc5, 0xc9, 0x19, 0x8, 0x96, 0xaf, 0x9b, 0xe6, 0xf6, 0x99, 0xb9, 0x58, 0x56, 0x5b, 0x25, 0xea, 0x9c, 0x86, 0x5e, 0x96, 0x6a, 0x48, 0xb},
		{0x75, 0xb2, 0xae, 0x1d, 0xd8, 0xca, 0xe6, 0x4d, 0xa8, 0xc5, 0xc9, 0x19, 0x8, 0x96, 0xaf, 0x9b, 0xe6, 0xf6, 0x99, 0xb9, 0x58, 0x56, 0x5b, 0x25, 0xea, 0x9c, 0x86, 0x5e, 0x96, 0x6a, 0x48, 0xb},
		{0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0x76, 0x88, 0xa0, 0x68, 0x45, 0x25, 0x8f, 0xd5, 0xf9, 0xb2, 0xb0, 0x42, 0x68, 0x6b, 0x51, 0xcc, 0x29, 0x94, 0x63, 0x85, 0xec, 0xf5, 0x47, 0xf0, 0x9c, 0x46, 0x86, 0xa9, 0x99, 0x7d, 0x29, 0x6c},
		{0x41, 0x44, 0x52, 0xff, 0x8c, 0xa6, 0xb3, 0x2e, 0xcc, 0x5e, 0x63, 0x8f, 0x8e, 0x7d, 0xe7, 0x52, 0x42, 0x94, 0x55, 0x2f, 0x89, 0xdd, 0x1e, 0x3c, 0xb0, 0xf4, 0x51, 0x51, 0x36, 0x81, 0x72, 0x1},
		{0xa9, 0xbb, 0x6a, 0x1f, 0x5d, 0x86, 0x7d, 0xa7, 0x5a, 0x7d, 0x9d, 0x8d, 0xc0, 0x15, 0xb7, 0x0, 0xee, 0xa9, 0x68, 0x51, 0x88, 0x57, 0x5a, 0xd9, 0x4e, 0x1d, 0x8e, 0x44, 0xbf, 0xdc, 0x73, 0xff},
		{0xf9, 0x11, 0x2c, 0xc2, 0x71, 0x70, 0xde, 0x47, 0x26, 0xeb, 0x26, 0xd4, 0xa4, 0xe8, 0x68, 0xb, 0x16, 0xa2, 0x6e, 0x52, 0x54, 0xe, 0x5c, 0x83, 0x17, 0x3, 0xea, 0xdd, 0xd5, 0xa7, 0xb2, 0x3f},
		{0x3d, 0xf3, 0x66, 0xd4, 0x12, 0x40, 0x3f, 0x28, 0xeb, 0xe4, 0x19, 0x59, 0xae, 0xab, 0x4d, 0xf3, 0x98, 0x88, 0x7f, 0x1e, 0x58, 0xa, 0x5d, 0xd4, 0xeb, 0xe5, 0x5d, 0x3d, 0x11, 0x70, 0x24, 0x76},
		{0xd6, 0x4c, 0xb1, 0xac, 0x61, 0x7, 0x26, 0xbb, 0xd3, 0x27, 0x2a, 0xcd, 0xdd, 0x55, 0xf, 0x2b, 0x6a, 0xe8, 0x1, 0x31, 0x48, 0x66, 0x2f, 0x98, 0x7b, 0x6d, 0x27, 0x69, 0xd9, 0x40, 0xcc, 0x37},
		{0xbc, 0xbb, 0x39, 0x57, 0x61, 0x1d, 0x54, 0xd6, 0x1b, 0xfe, 0x7a, 0xbd, 0x29, 0x52, 0x57, 0xdd, 0x19, 0x1, 0x89, 0x22, 0x7d, 0xdf, 0x7b, 0x53, 0x9f, 0xb, 0x46, 0x5, 0x9f, 0x80, 0xcc, 0x8e},
	}
	assert.DeepEqual(t, expected, root)
}

func TestComputeFieldRootsWithHasher_Capella(t *testing.T) {
	beaconState, err := util.NewBeaconStateCapella(util.FillRootsNaturalOptCapella)
	require.NoError(t, err)
	require.NoError(t, beaconState.SetGenesisTime(123))
	require.NoError(t, beaconState.SetGenesisValidatorsRoot(genesisValidatorsRoot()))
	require.NoError(t, beaconState.SetSlot(123))
	require.NoError(t, beaconState.SetFork(fork()))
	require.NoError(t, beaconState.SetLatestBlockHeader(latestBlockHeader()))
	historicalRoots, err := util.PrepareRoots(int(params.BeaconConfig().SlotsPerHistoricalRoot))
	require.NoError(t, err)
	require.NoError(t, beaconState.SetHistoricalRoots(historicalRoots))
	require.NoError(t, beaconState.SetEth1Data(eth1Data()))
	require.NoError(t, beaconState.SetEth1DataVotes([]*ethpb.Eth1Data{eth1Data()}))
	require.NoError(t, beaconState.SetEth1DepositIndex(123))
	require.NoError(t, beaconState.SetValidators([]*ethpb.Validator{validator()}))
	require.NoError(t, beaconState.SetBalances([]uint64{1, 2, 3}))
	randaoMixes, err := util.PrepareRoots(int(params.BeaconConfig().EpochsPerHistoricalVector))
	require.NoError(t, err)
	require.NoError(t, beaconState.SetRandaoMixes(randaoMixes))
	require.NoError(t, beaconState.SetSlashings([]uint64{1, 2, 3}))
	require.NoError(t, beaconState.SetPreviousParticipationBits([]byte{1, 2, 3}))
	require.NoError(t, beaconState.SetCurrentParticipationBits([]byte{1, 2, 3}))
	require.NoError(t, beaconState.SetJustificationBits(justificationBits()))
	require.NoError(t, beaconState.SetPreviousJustifiedCheckpoint(checkpoint("previous")))
	require.NoError(t, beaconState.SetCurrentJustifiedCheckpoint(checkpoint("current")))
	require.NoError(t, beaconState.SetFinalizedCheckpoint(checkpoint("finalized")))
	require.NoError(t, beaconState.SetInactivityScores([]uint64{1, 2, 3}))
	require.NoError(t, beaconState.SetCurrentSyncCommittee(syncCommittee("current")))
	require.NoError(t, beaconState.SetNextSyncCommittee(syncCommittee("next")))
	wrappedHeader, err := blocks.WrappedExecutionPayloadHeaderCapella(executionPayloadHeaderCapella(), big.NewInt(0))
	require.NoError(t, err)
	require.NoError(t, beaconState.SetLatestExecutionPayloadHeader(wrappedHeader))
	require.NoError(t, beaconState.SetNextWithdrawalIndex(123))
	require.NoError(t, beaconState.SetNextWithdrawalValidatorIndex(123))
	require.NoError(t, beaconState.AppendHistoricalSummaries(&ethpb.HistoricalSummary{
		BlockSummaryRoot: bytesutil.PadTo([]byte("block summary root"), 32),
		StateSummaryRoot: bytesutil.PadTo([]byte("state summary root"), 32),
	}))

	nativeState, ok := beaconState.(*statenative.BeaconState)
	require.Equal(t, true, ok)
	protoState, ok := nativeState.ToProtoUnsafe().(*ethpb.BeaconStateCapella)
	require.Equal(t, true, ok)
	initState, err := statenative.InitializeFromProtoCapella(protoState)
	require.NoError(t, err)
	s, ok := initState.(*statenative.BeaconState)
	require.Equal(t, true, ok)

	root, err := statenative.ComputeFieldRootsWithHasher(context.Background(), s)
	require.NoError(t, err)
	expected := [][]byte{
		{0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0x67, 0x76, 0x72, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0x58, 0xba, 0xf, 0x9b, 0x4f, 0x63, 0x1c, 0xa6, 0x19, 0xb1, 0xa2, 0x1f, 0xd1, 0x29, 0xc7, 0x67, 0x9c, 0x32, 0x4, 0x1f, 0xcf, 0x4e, 0x64, 0x9b, 0x8f, 0x21, 0xb4, 0xe6, 0xa5, 0xc9, 0xc, 0x38},
		{0x8b, 0x5, 0x59, 0x78, 0xed, 0xbe, 0x2c, 0xde, 0xa6, 0xf, 0x52, 0xdc, 0x16, 0x83, 0xa0, 0x5d, 0x8, 0xc3, 0x37, 0x91, 0x3a, 0xf6, 0xfa, 0x6, 0x62, 0xc9, 0x6, 0xb1, 0x41, 0x48, 0xaf, 0xec},
		{0xcc, 0xa1, 0x69, 0xa8, 0x8b, 0x1e, 0x49, 0x5, 0x1f, 0xe4, 0xd8, 0xff, 0x82, 0xa2, 0x2e, 0xf0, 0x54, 0xd1, 0x13, 0xc9, 0x8e, 0xb9, 0x82, 0xa6, 0x9e, 0x42, 0x2, 0xec, 0x97, 0x6f, 0x33, 0x88},
		{0xcc, 0xa1, 0x69, 0xa8, 0x8b, 0x1e, 0x49, 0x5, 0x1f, 0xe4, 0xd8, 0xff, 0x82, 0xa2, 0x2e, 0xf0, 0x54, 0xd1, 0x13, 0xc9, 0x8e, 0xb9, 0x82, 0xa6, 0x9e, 0x42, 0x2, 0xec, 0x97, 0x6f, 0x33, 0x88},
		{0xf, 0xad, 0xd3, 0x92, 0x4b, 0xda, 0xfa, 0xc6, 0x61, 0x50, 0xb7, 0xdf, 0x5f, 0x2c, 0xd0, 0x94, 0xc3, 0xaf, 0x41, 0x9d, 0xa, 0xea, 0x50, 0x96, 0x82, 0x62, 0x1c, 0x72, 0x26, 0x20, 0x6b, 0xac},
		{0xc9, 0x4e, 0x2c, 0xb0, 0x20, 0xe3, 0xe7, 0x8c, 0x5c, 0xbd, 0xeb, 0x9b, 0xa5, 0x7b, 0x53, 0x50, 0xca, 0xfe, 0xe9, 0x48, 0x9e, 0x8d, 0xf8, 0x4a, 0xe6, 0x8d, 0x9c, 0x97, 0x81, 0x74, 0xb, 0x5e},
		{0x71, 0x52, 0xd2, 0x9b, 0x87, 0x3c, 0x8a, 0xd9, 0x51, 0x55, 0xc0, 0x42, 0xb, 0xc4, 0x12, 0xa4, 0x79, 0xf5, 0x7d, 0x37, 0x16, 0xf4, 0x90, 0x72, 0x5d, 0xe0, 0x34, 0xb4, 0x2, 0x8c, 0x39, 0xe4},
		{0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0xf, 0xf7, 0x4f, 0xe1, 0xa9, 0x72, 0x9c, 0x95, 0xf0, 0xe1, 0xde, 0xa4, 0x32, 0xc, 0x67, 0x52, 0x23, 0x13, 0x9e, 0xe2, 0x40, 0x8d, 0xf6, 0x18, 0x57, 0xf0, 0x1a, 0x4a, 0xad, 0x46, 0xce, 0x42},
		{0xf9, 0x11, 0x2c, 0xc2, 0x71, 0x70, 0xde, 0x47, 0x26, 0xeb, 0x26, 0xd4, 0xa4, 0xe8, 0x68, 0xb, 0x16, 0xa2, 0x6e, 0x52, 0x54, 0xe, 0x5c, 0x83, 0x17, 0x3, 0xea, 0xdd, 0xd5, 0xa7, 0xb2, 0x3f},
		{0x64, 0xbd, 0x40, 0xa7, 0x10, 0x44, 0x84, 0xed, 0xf3, 0x5f, 0xc3, 0x5d, 0x7b, 0xbe, 0xe8, 0x75, 0xbf, 0x66, 0xcb, 0xce, 0x77, 0xfa, 0x0, 0x3, 0xdd, 0xfb, 0x80, 0xd2, 0x77, 0x1b, 0xc2, 0x8},
		{0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x3, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0x75, 0xb2, 0xae, 0x1d, 0xd8, 0xca, 0xe6, 0x4d, 0xa8, 0xc5, 0xc9, 0x19, 0x8, 0x96, 0xaf, 0x9b, 0xe6, 0xf6, 0x99, 0xb9, 0x58, 0x56, 0x5b, 0x25, 0xea, 0x9c, 0x86, 0x5e, 0x96, 0x6a, 0x48, 0xb},
		{0x75, 0xb2, 0xae, 0x1d, 0xd8, 0xca, 0xe6, 0x4d, 0xa8, 0xc5, 0xc9, 0x19, 0x8, 0x96, 0xaf, 0x9b, 0xe6, 0xf6, 0x99, 0xb9, 0x58, 0x56, 0x5b, 0x25, 0xea, 0x9c, 0x86, 0x5e, 0x96, 0x6a, 0x48, 0xb},
		{0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0x76, 0x88, 0xa0, 0x68, 0x45, 0x25, 0x8f, 0xd5, 0xf9, 0xb2, 0xb0, 0x42, 0x68, 0x6b, 0x51, 0xcc, 0x29, 0x94, 0x63, 0x85, 0xec, 0xf5, 0x47, 0xf0, 0x9c, 0x46, 0x86, 0xa9, 0x99, 0x7d, 0x29, 0x6c},
		{0x41, 0x44, 0x52, 0xff, 0x8c, 0xa6, 0xb3, 0x2e, 0xcc, 0x5e, 0x63, 0x8f, 0x8e, 0x7d, 0xe7, 0x52, 0x42, 0x94, 0x55, 0x2f, 0x89, 0xdd, 0x1e, 0x3c, 0xb0, 0xf4, 0x51, 0x51, 0x36, 0x81, 0x72, 0x1},
		{0xa9, 0xbb, 0x6a, 0x1f, 0x5d, 0x86, 0x7d, 0xa7, 0x5a, 0x7d, 0x9d, 0x8d, 0xc0, 0x15, 0xb7, 0x0, 0xee, 0xa9, 0x68, 0x51, 0x88, 0x57, 0x5a, 0xd9, 0x4e, 0x1d, 0x8e, 0x44, 0xbf, 0xdc, 0x73, 0xff},
		{0xf9, 0x11, 0x2c, 0xc2, 0x71, 0x70, 0xde, 0x47, 0x26, 0xeb, 0x26, 0xd4, 0xa4, 0xe8, 0x68, 0xb, 0x16, 0xa2, 0x6e, 0x52, 0x54, 0xe, 0x5c, 0x83, 0x17, 0x3, 0xea, 0xdd, 0xd5, 0xa7, 0xb2, 0x3f},
		{0x3d, 0xf3, 0x66, 0xd4, 0x12, 0x40, 0x3f, 0x28, 0xeb, 0xe4, 0x19, 0x59, 0xae, 0xab, 0x4d, 0xf3, 0x98, 0x88, 0x7f, 0x1e, 0x58, 0xa, 0x5d, 0xd4, 0xeb, 0xe5, 0x5d, 0x3d, 0x11, 0x70, 0x24, 0x76},
		{0xd6, 0x4c, 0xb1, 0xac, 0x61, 0x7, 0x26, 0xbb, 0xd3, 0x27, 0x2a, 0xcd, 0xdd, 0x55, 0xf, 0x2b, 0x6a, 0xe8, 0x1, 0x31, 0x48, 0x66, 0x2f, 0x98, 0x7b, 0x6d, 0x27, 0x69, 0xd9, 0x40, 0xcc, 0x37},
		{0x39, 0x29, 0x16, 0xe8, 0x5a, 0xd2, 0xb, 0xbb, 0x1f, 0xef, 0x6a, 0xe0, 0x2d, 0xa6, 0x6a, 0x46, 0x81, 0xba, 0xcf, 0x86, 0xfc, 0x16, 0x22, 0x2a, 0x9b, 0x72, 0x96, 0x71, 0x2b, 0xc7, 0x5b, 0x9d},
		{0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0x7b, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		{0xa1, 0x4, 0x64, 0x31, 0x2a, 0xa, 0x49, 0x31, 0x1c, 0x1, 0x41, 0x17, 0xc0, 0x52, 0x52, 0xfa, 0x4c, 0xf4, 0x95, 0x4f, 0x5c, 0xb0, 0x5a, 0x40, 0xc1, 0x32, 0x39, 0xc3, 0x7c, 0xb7, 0x2c, 0x27},
	}
	assert.DeepEqual(t, expected, root)
}

func genesisValidatorsRoot() []byte {
	gvr := bytesutil.ToBytes32([]byte("gvr"))
	return gvr[:]
}

func fork() *ethpb.Fork {
	prev := bytesutil.ToBytes4([]byte("prev"))
	curr := bytesutil.ToBytes4([]byte("curr"))
	return &ethpb.Fork{
		PreviousVersion: prev[:],
		CurrentVersion:  curr[:],
		Epoch:           123,
	}
}

func latestBlockHeader() *ethpb.BeaconBlockHeader {
	pr := bytesutil.ToBytes32([]byte("parent"))
	sr := bytesutil.ToBytes32([]byte("state"))
	br := bytesutil.ToBytes32([]byte("body"))
	return &ethpb.BeaconBlockHeader{
		Slot:          123,
		ProposerIndex: 123,
		ParentRoot:    pr[:],
		StateRoot:     sr[:],
		BodyRoot:      br[:],
	}
}

func eth1Data() *ethpb.Eth1Data {
	dr := bytesutil.ToBytes32([]byte("deposit"))
	bh := bytesutil.ToBytes32([]byte("block"))
	return &ethpb.Eth1Data{
		DepositRoot:  dr[:],
		DepositCount: 123,
		BlockHash:    bh[:],
	}
}

func validator() *ethpb.Validator {
	pk := bytesutil.ToBytes48([]byte("public"))
	wc := bytesutil.ToBytes32([]byte("withdrawal"))
	return &ethpb.Validator{
		PublicKey:                  pk[:],
		WithdrawalCredentials:      wc[:],
		EffectiveBalance:           123,
		Slashed:                    true,
		ActivationEligibilityEpoch: 123,
		ActivationEpoch:            123,
		ExitEpoch:                  123,
		WithdrawableEpoch:          123,
	}
}

func pendingAttestation(prefix string) *ethpb.PendingAttestation {
	bbr := bytesutil.ToBytes32([]byte(prefix + "beacon"))
	r := bytesutil.ToBytes32([]byte(prefix + "root"))
	return &ethpb.PendingAttestation{
		AggregationBits: bitfield.Bitlist{0x00, 0xFF, 0xFF, 0xFF},
		Data: &ethpb.AttestationData{
			Slot:            123,
			CommitteeIndex:  123,
			BeaconBlockRoot: bbr[:],
			Source: &ethpb.Checkpoint{
				Epoch: 123,
				Root:  r[:],
			},
			Target: &ethpb.Checkpoint{
				Epoch: 123,
				Root:  r[:],
			},
		},
		InclusionDelay: 123,
		ProposerIndex:  123,
	}
}

func justificationBits() bitfield.Bitvector4 {
	v := bitfield.NewBitvector4()
	v.SetBitAt(1, true)
	return v
}

func checkpoint(prefix string) *ethpb.Checkpoint {
	r := bytesutil.ToBytes32([]byte(prefix + "root"))
	return &ethpb.Checkpoint{
		Epoch: 123,
		Root:  r[:],
	}
}

func syncCommittee(prefix string) *ethpb.SyncCommittee {
	pubkeys := make([][]byte, params.BeaconConfig().SyncCommitteeSize)
	for i := range pubkeys {
		key := bytesutil.ToBytes48([]byte(prefix + "pubkey"))
		pubkeys[i] = key[:]
	}
	agg := bytesutil.ToBytes48([]byte(prefix + "aggregate"))
	return &ethpb.SyncCommittee{
		Pubkeys:         pubkeys,
		AggregatePubkey: agg[:],
	}
}

func executionPayloadHeader() *enginev1.ExecutionPayloadHeader {
	ph := bytesutil.ToBytes32([]byte("parent"))
	fr := bytesutil.PadTo([]byte("fee"), 20)
	sr := bytesutil.ToBytes32([]byte("state"))
	rr := bytesutil.ToBytes32([]byte("receipts"))
	lb := bytesutil.PadTo([]byte("logs"), 256)
	pr := bytesutil.ToBytes32([]byte("prev"))
	ed := bytesutil.ToBytes32([]byte("extra"))
	bf := bytesutil.ToBytes32([]byte("base"))
	bh := bytesutil.ToBytes32([]byte("block"))
	tr := bytesutil.ToBytes32([]byte("transactions"))
	return &enginev1.ExecutionPayloadHeader{
		ParentHash:       ph[:],
		FeeRecipient:     fr,
		StateRoot:        sr[:],
		ReceiptsRoot:     rr[:],
		LogsBloom:        lb,
		PrevRandao:       pr[:],
		BlockNumber:      0,
		GasLimit:         0,
		GasUsed:          0,
		Timestamp:        0,
		ExtraData:        ed[:],
		BaseFeePerGas:    bf[:],
		BlockHash:        bh[:],
		TransactionsRoot: tr[:],
	}
}

func executionPayloadHeaderCapella() *enginev1.ExecutionPayloadHeaderCapella {
	ph := bytesutil.ToBytes32([]byte("parent"))
	fr := bytesutil.PadTo([]byte("fee"), 20)
	sr := bytesutil.ToBytes32([]byte("state"))
	rr := bytesutil.ToBytes32([]byte("receipts"))
	lb := bytesutil.PadTo([]byte("logs"), 256)
	pr := bytesutil.ToBytes32([]byte("prev"))
	ed := bytesutil.ToBytes32([]byte("extra"))
	bf := bytesutil.ToBytes32([]byte("base"))
	bh := bytesutil.ToBytes32([]byte("block"))
	tr := bytesutil.ToBytes32([]byte("transactions"))
	wr := bytesutil.ToBytes32([]byte("withdrawals"))
	return &enginev1.ExecutionPayloadHeaderCapella{
		ParentHash:       ph[:],
		FeeRecipient:     fr,
		StateRoot:        sr[:],
		ReceiptsRoot:     rr[:],
		LogsBloom:        lb,
		PrevRandao:       pr[:],
		BlockNumber:      0,
		GasLimit:         0,
		GasUsed:          0,
		Timestamp:        0,
		ExtraData:        ed[:],
		BaseFeePerGas:    bf[:],
		BlockHash:        bh[:],
		TransactionsRoot: tr[:],
		WithdrawalsRoot:  wr[:],
	}
}
