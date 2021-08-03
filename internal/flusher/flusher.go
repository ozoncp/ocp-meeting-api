package flusher

import (
	"github.com/ozoncp/ocp-meeting-api/internal/models"
	"github.com/ozoncp/ocp-meeting-api/internal/repo"
	"github.com/ozoncp/ocp-meeting-api/internal/utils"
)

type Flusher interface {
	Flush(tasks []models.Meeting) []models.Meeting
}

func NewFlusher(chunkSize int, meetingRepo repo.Repo) Flusher {
	return &flusher{
		chunkSize:   chunkSize,
		meetingRepo: meetingRepo,
	}
}

type flusher struct {
	chunkSize   int
	meetingRepo repo.Repo
}

func (f *flusher) Flush(meetings []models.Meeting) []models.Meeting {
	chunks := utils.SplitToBulks(meetings, uint(f.chunkSize))

	for i, meeting := range chunks {
		if err := f.meetingRepo.Add(meeting); err != nil {
			return meeting[i:]
		}
	}

	return nil
}
