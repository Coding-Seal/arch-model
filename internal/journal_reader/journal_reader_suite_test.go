package journalreader_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestJournalReader(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "JournalReader Suite")
}
