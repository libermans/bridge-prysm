package fork_transition

import (
	"testing"

	"github.com/OffchainLabs/prysm/v6/testing/spectest/shared/electra/fork"
)

func TestMainnet_Electra_Transition(t *testing.T) {
	fork.RunForkTransitionTest(t, "mainnet")
}
