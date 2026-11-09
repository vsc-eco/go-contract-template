package test

import (
	"fmt"
	"os"
	"testing"
)

// TestMain ensures the libp2p test config binds to a random port. The upstream
// default is 10720, which collides with anything already listening (e.g. a
// Docker container forwarding the same port on the host, or a second contract
// test suite running in parallel). Writing data/config/p2pConfig.json before
// tests run causes test_utils.NewContractTest to load Port=0 and let the kernel
// pick a free port.
func TestMain(m *testing.M) {
	if err := os.MkdirAll("data/config", 0755); err != nil {
		panic(fmt.Sprintf("failed to create data/config: %v", err))
	}
	cfg := []byte(`{"Port":0,"ServerMode":false,"AllowPrivate":false,"Bootnodes":[],"AnnounceAddrs":[]}`)
	if err := os.WriteFile("data/config/p2pConfig.json", cfg, 0644); err != nil {
		panic(fmt.Sprintf("failed to write p2pConfig.json: %v", err))
	}
	os.Exit(m.Run())
}
