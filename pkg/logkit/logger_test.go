package logkit_test

import (
	"github.com/airdb/chat-gateway/pkg/logkit"
	. "github.com/onsi/ginkgo/v2"
	"golang.org/x/exp/slog"
)

var _ = Describe("LogKit Init", func() {
	BeforeEach(func() {
		logkit.Init(logkit.WrapOptions(logkit.WithLevel(slog.LevelInfo)))
	})
	Describe("Getting External login url", func() {
		// positive test case
		Context("Get an error if returned url is empty", func() {
			It("Failed precondition err must be responded", func() {
				logkit.Log.Info("this is an info log message")
				// Expect(err).Should(BeNil())
				// Expect(res).ShouldNot(BeNil())
			})
		})
	})
})
