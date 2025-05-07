package ssz_static

import (
	"testing"

	"github.com/OffchainLabs/prysm/v6/testing/spectest/shared/capella/ssz_static"
)

func TestMinimal_Capella_SSZStatic(t *testing.T) {
	ssz_static.RunSSZStaticTests(t, "minimal")
}
