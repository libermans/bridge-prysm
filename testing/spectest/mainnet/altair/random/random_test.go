package random

import (
	"testing"

	"github.com/OffchainLabs/prysm/v6/testing/spectest/shared/altair/sanity"
)

func TestMainnet_Altair_Random(t *testing.T) {
	sanity.RunBlockProcessingTest(t, "mainnet", "random/random/pyspec_tests")
}
