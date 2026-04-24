package docker_test

import (
	"testing"

	"shipyard/internal/docker"
)

func TestNewClient(t *testing.T) {
	// Verifies the client can be constructed (requires Docker daemon running).
	// In CI without Docker, this test is skipped.
	c, err := docker.NewClient("shipyard_net")
	if err != nil {
		t.Skipf("docker not available: %v", err)
	}
	if c == nil {
		t.Error("expected non-nil client")
	}
}
