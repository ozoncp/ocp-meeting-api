package saver

import (
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/ozoncp/ocp-meeting-api/internal/mocks"
	"github.com/ozoncp/ocp-meeting-api/internal/models"
	"time"
)

var _ = Describe("Saver", func() {

	var (
		ctrl        *gomock.Controller
		mockFlusher *mocks.MockFlusher
		meetings    []models.Meeting
		meetingChan chan models.Meeting
	)

	var now = time.Now()

	BeforeEach(func() {

		ctrl = gomock.NewController(GinkgoT())
		mockFlusher = mocks.NewMockFlusher(ctrl)
		meetingChan = make(chan models.Meeting, 5)
		meetings = []models.Meeting{
			{1, 1, "", now, now},
			{2, 2, "", now, now},
			{3, 3, "", now, now},
			{4, 4, "", now, now},
			{5, 5, "", now, now},
		}

	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("Tests", func() {
		BeforeEach(func() {
			mockFlusher.EXPECT().Flush(gomock.Any()).MinTimes(1).Return(meetings)
		})

		It("Saving", func() {
			s := NewSaver(5, mockFlusher, time.Second)
			s.Init()

			for _, meeting := range meetings {
				_ = s.Save(meeting)
				meetingChan <- meeting
			}

			s.Close()
			Expect(len(meetingChan)).Should(BeEquivalentTo(5))
		})

		It("Panic", func() {
			mockFlusher.Flush(nil)
			s := NewSaver(5, mockFlusher, time.Second)
			s.Init()

			save := func() {
				_ = s.Save(meetings[0])
			}

			s.Close()
			Expect(save).Should(Panic())
		})
	})

})
