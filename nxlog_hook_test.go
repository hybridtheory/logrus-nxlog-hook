package nxlog

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"reflect"
	"time"
)

var _ = Describe("Nxlog hook", func() {

	var (
		entry *logrus.Entry
	)

	BeforeEach(func() {
		logger := logrus.New()
		logger.Out = &bytes.Buffer{}
		entry = &logrus.Entry{
			Time:    time.Time{},
			Level:   logrus.InfoLevel,
			Message: "test message",
			Data:    logrus.Fields{},
			Logger:  logger,
		}
	})

	It("can be instantiated", func() {
		hook, err := NewNxlogHook("tcp", "127.0.0.1:3000", map[string]interface{}{})
		Expect(err).To(BeNil())
		Expect(reflect.TypeOf(hook).String()).To(Equal("*nxlog.Hook"))
	})

	It("returns an error if connection fails", func() {
		hook, err := NewNxlogHook("tcp", "127.0.0.1:7005", map[string]interface{}{})
		Expect(err).ToNot(BeNil())
		Expect(hook).To(BeNil())
	})

	It("allows the setting the minimum log level that will trigger this hook", func() {
		hook, _ := NewNxlogHook("tcp", "127.0.0.1:3000", map[string]interface{}{})
		hook.Level = logrus.FatalLevel
		Expect(hook.Levels()).ToNot(ContainElement(logrus.InfoLevel))
		Expect(hook.Levels()).ToNot(ContainElement(logrus.ErrorLevel))
		hook.Level = logrus.ErrorLevel
		Expect(hook.Levels()).ToNot(ContainElement(logrus.InfoLevel))
		Expect(hook.Levels()).To(ContainElement(logrus.ErrorLevel))
		hook.Level = logrus.DebugLevel
		Expect(hook.Levels()).To(ContainElement(logrus.InfoLevel))
		Expect(hook.Levels()).To(ContainElement(logrus.ErrorLevel))
	})

	It("allows formatting a message", func() {
		hook, _ := NewNxlogHook("tcp", "127.0.0.1:3000", map[string]interface{}{})
		By("using the default formatter", func() {
			formatted, _ := hook.getMessage(entry)
			Expect(formatted).To(Equal(`{"level":"info","msg":"test message","time":"0001-01-01T00:00:00Z"}` + "\n"))
		})

		By("using a text formatter", func() {
			hook.Formatter = &logrus.TextFormatter{}
			formatted, _ := hook.getMessage(entry)
			Expect(formatted).To(Equal(`time="0001-01-01T00:00:00Z" level=info msg="test message"` + "\n"))
		})
	})
})
