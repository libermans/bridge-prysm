package operations

import (
	"testing"

	"github.com/OffchainLabs/prysm/v6/testing/spectest/shared/bellatrix/operations"
)

func TestMinimal_Bellatrix_Operations_Attestation(t *testing.T) {
	operations.RunAttestationTest(t, "minimal")
}
