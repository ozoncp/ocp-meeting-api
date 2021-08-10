package saver

import (
	"github.com/ozoncp/ocp-meeting-api/internal/flusher"
	"github.com/ozoncp/ocp-meeting-api/internal/models"
	"time"
)

type Saver interface {
	Save(meeting models.Meeting) error
	Init()
	Close()
}

type saver struct {
	flusher     flusher.Flusher
	meetingChan chan models.Meeting
	closeChan   chan struct{}
	duration    time.Duration
}

// NewSaver возвращает Saver с поддержкой переодического сохранения
func NewSaver(capacity uint, flusher flusher.Flusher, duration time.Duration) Saver {
	return &saver{
		flusher:     flusher,
		meetingChan: make(chan models.Meeting, capacity),
		closeChan:   make(chan struct{}),
		duration:    duration,
	}
}

func (s *saver) Save(meeting models.Meeting) error {
	s.meetingChan <- meeting
	return nil
}

func (s *saver) Init() {
	go func() {
		meetings := make([]models.Meeting, 0)

		ticker := time.NewTicker(s.duration)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				_ = s.flusher.Flush(meetings)
			case meeting := <-s.meetingChan:
				meetings = append(meetings, meeting)
			case <-s.closeChan:
				_ = s.flusher.Flush(meetings)
				close(s.meetingChan)
				return
			}
		}
	}()
}

func (s *saver) Close() {
	s.closeChan <- struct{}{}
}
