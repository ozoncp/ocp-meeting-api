package flusher

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ozoncp/ocp-meeting-api/internal/mocks"
	"github.com/ozoncp/ocp-meeting-api/internal/models"
	"time"
)

var _ = Describe("Flusher", func() {
	var (
		ctrl     *gomock.Controller
		mockRepo *mocks.MockRepo
		f        Flusher
		meetings []models.Meeting
		result   []models.Meeting
		ctx      context.Context
	)

	var now = time.Now()

	BeforeEach(func() {

		ctrl = gomock.NewController(GinkgoT())
		mockRepo = mocks.NewMockRepo(ctrl)
		ctx = context.Background()

		meetings = []models.Meeting{
			{1, 1, "", now, now},
			{2, 2, "", now, now},
			{3, 3, "", now, now},
			{4, 4, "", now, now},
			{5, 5, "", now, now},
		}

	})

	JustBeforeEach(func() {
		f = NewFlusher(3, mockRepo)
		result = f.Flush(ctx, meetings)
	})

	Context("Save all", func() {
		BeforeEach(func() {
			mockRepo.EXPECT().AddMany(ctx, gomock.Any()).Return(nil).AnyTimes()
		})

		It("", func() {
			Expect(result).Should(BeNil())
		})
	})

	Context("Saving error", func() {
		BeforeEach(func() {
			mockRepo.EXPECT().AddMany(ctx, gomock.Any()).Return(errors.New("error"))
		})

		It("", func() {
			Expect(len(result)).Should(BeEquivalentTo(len(meetings)))
			Expect(result).Should(BeEquivalentTo(meetings))
		})
	})

	Context("Partial saving", func() {
		BeforeEach(func() {
			mockRepo.EXPECT().AddMany(ctx, gomock.Any()).Return(errors.New("error"))
			mockRepo.EXPECT().AddMany(ctx, gomock.Any()).Return(nil).Times(1)
		})

		It("", func() {
			Expect(result).Should(BeEquivalentTo(meetings))
		})
	})
})
