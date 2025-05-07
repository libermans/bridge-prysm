package operations

import (
	"testing"

	"github.com/OffchainLabs/prysm/v6/testing/spectest/shared/electra/operations"
)

func TestMinimal_Electra_Operations_BlockHeader(t *testing.T) {
	operations.RunBlockHeaderTest(t, "minimal")
}
