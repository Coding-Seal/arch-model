package dequeue_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDequeue(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dequeue Suite")
}
