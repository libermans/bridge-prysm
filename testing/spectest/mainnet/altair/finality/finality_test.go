package finality

import (
	"testing"

	"github.com/OffchainLabs/prysm/v6/testing/spectest/shared/altair/finality"
)

func TestMainnet_Altair_Finality(t *testing.T) {
	finality.RunFinalityTest(t, "mainnet")
}
