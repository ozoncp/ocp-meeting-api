package repo

import "github.com/ozoncp/ocp-meeting-api/internal/models"

type Repo interface {
	Add(meetings []models.Meeting) error
	Describe(meetingId uint64) (*models.Meeting, error)
	List(limit, offset uint64) ([]models.Meeting, error)
}
