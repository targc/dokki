package basic

import (
	"testing"

	"github.com/targc/dokki"
)

func TestExampleBasic(t *testing.T) {
	s := dokki.NewSetup()

	defer s.Down(t)

	_ = s.SetupTestSuiteDockerCompose(
		t,
		"docker-compose.yml",
		[]string{
			"redis",
		},
	)
}
