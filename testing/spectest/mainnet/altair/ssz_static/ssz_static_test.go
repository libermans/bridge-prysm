package ssz_static

import (
	"testing"

	"github.com/OffchainLabs/prysm/v6/testing/spectest/shared/altair/ssz_static"
)

func TestMainnet_Altair_SSZStatic(t *testing.T) {
	ssz_static.RunSSZStaticTests(t, "mainnet")
}
