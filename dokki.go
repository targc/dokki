package dokki

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go/modules/compose"
)

type TestSuite struct {
	composeStack compose.ComposeStack
	services     []string
}

func SetupTestSuite(t *testing.T, composeFile string, services []string) *TestSuite {
	composeStack, err := compose.NewDockerCompose(composeFile)

	if err != nil {
		t.Fatalf("Failed to create compose stack: %v", err)
	}

	ts := TestSuite{
		composeStack: composeStack,
		services:     services,
	}

	t.Cleanup(func() {
		ts.teardown(t)
	})

	ts.up(t)

	return &ts
}

func (ts *TestSuite) up(t *testing.T) {
	ctx := context.Background()

	err := ts.composeStack.Up(
		ctx,
		compose.Wait(true),
		compose.RunServices(
			ts.services...,
		),
	)

	if err != nil {
		t.Fatalf("Failed to start compose stack: %v", err)
	}
}

func (ts *TestSuite) teardown(t *testing.T) {
	if ts.composeStack == nil {
		return
	}

	ctx := context.Background()

	err := ts.composeStack.Down(
		ctx,
		compose.RemoveOrphans(true),
		compose.RemoveVolumes(true),
		compose.RemoveImagesLocal,
	)

	if err != nil {
		t.Fatalf("Failed to down compose stack: %v", err)
	}
}
