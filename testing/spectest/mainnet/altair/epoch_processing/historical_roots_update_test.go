package epoch_processing

import (
	"testing"

	"github.com/OffchainLabs/prysm/v6/testing/spectest/shared/altair/epoch_processing"
)

func TestMainnet_Altair_EpochProcessing_HistoricalRootsUpdate(t *testing.T) {
	epoch_processing.RunHistoricalRootsUpdateTests(t, "mainnet")
}
