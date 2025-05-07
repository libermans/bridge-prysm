package epoch_processing

import (
	"testing"

	"github.com/OffchainLabs/prysm/v6/testing/spectest/shared/capella/epoch_processing"
)

func TestMainnet_Capella_EpochProcessing_Slashings(t *testing.T) {
	epoch_processing.RunSlashingsTests(t, "mainnet")
}
