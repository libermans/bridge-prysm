package operations

import (
	"testing"

	"github.com/OffchainLabs/prysm/v6/testing/spectest/shared/bellatrix/operations"
)

func TestMinimal_Bellatrix_Operations_Deposit(t *testing.T) {
	operations.RunDepositTest(t, "minimal")
}
