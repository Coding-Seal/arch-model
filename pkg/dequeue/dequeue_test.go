package dequeue_test

import (
	"github.com/Coding-Seal/arch-model/pkg/dequeue"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const (
	queueSize = 5
	testElem  = 69
)

var _ = Describe("Dequeue", func() {
	var q *dequeue.Dequeue[int]
	Context("Queue is empty", func() {
		BeforeEach(func() {
			q = dequeue.New[int](queueSize)
			Expect(q.Empty()).To(BeTrue())
		})
		It("Back", func() {
			elem, ok := q.Back()
			Expect(ok).To(BeFalse())
			Expect(elem).To(Equal(0))
		})
		It("Front", func() {
			elem, ok := q.Front()
			Expect(ok).To(BeFalse())
			Expect(elem).To(Equal(0))
		})
		It("PushFront", func() {
			Expect(q.PushFront(testElem)).To(BeTrue())

			By("front and back are equal to element and not empty")
			elem, ok := q.Front()
			Expect(ok).To(BeTrue())
			Expect(elem).To(Equal(testElem))

			elem, ok = q.Back()
			Expect(ok).To(BeTrue())
			Expect(elem).To(Equal(testElem))
		})
		It("PushBack", func() {
			Expect(q.PushBack(testElem)).To(BeTrue())

			By("front and back are equal to element and not empty")
			elem, ok := q.Front()
			Expect(ok).To(BeTrue())
			Expect(elem).To(Equal(testElem))

			elem, ok = q.Back()
			Expect(ok).To(BeTrue())
			Expect(elem).To(Equal(testElem))
		})
		It("PopFront", func() {
			Expect(q.PopFront()).To(BeFalse())
		})
		It("PopBack", func() {
			Expect(q.PopBack()).To(BeFalse())
		})
	})
	Context("Queue is full", func() {
		BeforeEach(func() {
			q = dequeue.New[int](queueSize)
			for i := range queueSize {
				Expect(q.PushBack(i + 1)).To(BeTrue())
			}

			Expect(q.Full()).To(BeTrue())
		})
		It("Back", func() {
			elem, ok := q.Back()
			Expect(ok).To(BeTrue())
			Expect(elem).To(Equal(queueSize))
		})
		It("Front", func() {
			elem, ok := q.Front()
			Expect(ok).To(BeTrue())
			Expect(elem).To(Equal(1))
		})
		It("PushFront", func() {
			Expect(q.PushFront(testElem)).To(BeFalse())
		})
		It("PushBack", func() {
			Expect(q.PushBack(testElem)).To(BeFalse())
		})
		It("PopFront", func() {
			Expect(q.PopFront()).To(BeTrue())
			elem, ok := q.Front()
			Expect(ok).To(BeTrue())
			Expect(elem).To(Equal(2))
		})
		It("PopBack", func() {
			Expect(q.PopBack()).To(BeTrue())
			elem, ok := q.Back()
			Expect(ok).To(BeTrue())
			Expect(elem).To(Equal(queueSize - 1))
		})
	})
	Context("Just some casual use", func() {
		BeforeEach(func() {
			q = dequeue.New[int](queueSize)
		})
		It("Use as Queue", func() {
			for i := range queueSize {
				Expect(q.PushBack(i + 1)).To(BeTrue())
			}

			Expect(q.Full()).To(BeTrue())

			for i := range queueSize {
				elem, ok := q.Front()
				Expect(ok).To(BeTrue())
				Expect(elem).To(Equal(i + 1))
				Expect(q.PopFront()).To(BeTrue())
			}
			Expect(q.Empty()).To(BeTrue())
		})
	})
})
