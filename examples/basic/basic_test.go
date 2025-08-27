package basic

import (
	"testing"

	"github.com/targc/dokki"
)

func TestExampleBasic(t *testing.T) {
	_ = dokki.SetupTestSuite(
		t,
		"docker-compose.yml",
		[]string{
			"redis",
		},
	)
}
