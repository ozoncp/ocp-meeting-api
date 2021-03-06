package flusher

import (
	"context"
	"github.com/ozoncp/ocp-meeting-api/internal/models"
	"github.com/ozoncp/ocp-meeting-api/internal/repo"
	"github.com/ozoncp/ocp-meeting-api/internal/utils"
)

type Flusher interface {
	Flush(ctx context.Context, meetings []models.Meeting) []models.Meeting
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

func (f *flusher) Flush(ctx context.Context, meetings []models.Meeting) []models.Meeting {
	chunks := utils.SplitToBulks(meetings, uint(f.chunkSize))
	var problemMeetings [][]models.Meeting

	for i, meeting := range chunks {
		if _, err := f.meetingRepo.AddMany(ctx, meeting); err != nil {
			problemMeetings = chunks[i:]
			break
		}
	}

	if problemMeetings != nil {
		problemList := make([]models.Meeting, 0, len(problemMeetings)*f.chunkSize)
		for _, meeting := range chunks {
			problemList = append(problemList, meeting...)
		}
		return problemList
	}

	return nil
}
