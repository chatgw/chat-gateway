package logkit_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLogkit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Logkit Suite")
}
