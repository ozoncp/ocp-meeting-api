package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/opentracing/opentracing-go"
	"github.com/ozoncp/ocp-meeting-api/internal/metrics"
	"github.com/ozoncp/ocp-meeting-api/internal/models"
	"github.com/ozoncp/ocp-meeting-api/internal/producer"
	"github.com/ozoncp/ocp-meeting-api/internal/repo"
	"github.com/ozoncp/ocp-meeting-api/internal/utils"
	desc "github.com/ozoncp/ocp-meeting-api/pkg/ocp-meeting-api"
	log "github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

const (
	batchSize uint = 5
)

type api struct {
	desc.UnimplementedOcpMeetingApiServer
	repo repo.Repo
	prod producer.Producer
}

func NewOcpMeetingApi(repo repo.Repo, prod producer.Producer) desc.OcpMeetingApiServer {
	return &api{
		repo: repo,
		prod: prod,
	}
}

func (a *api) MultiCreateMeetingsV1(
	ctx context.Context,
	req *desc.MultiCreateMeetingsV1Request,
) (*desc.MultiCreateMeetingsV1Response, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	var meetings []models.Meeting
	for _, meeting := range req.Meetings {
		meetings = append(meetings, models.Meeting{
			Id:     meeting.Id,
			UserId: meeting.UserId,
			Link:   meeting.Link,
			Start:  meeting.Start.AsTime(),
			End:    meeting.End.AsTime(),
		})
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "MultiCreateMeetings")
	defer span.Finish()

	bulks := utils.SplitToBulks(meetings, batchSize)
	response := &desc.MultiCreateMeetingsV1Response{}

	for i := 0; i < len(bulks); i++ {
		meetingIds, err := a.repo.AddMany(ctx, bulks[i])
		if err != nil {
			log.Error().Msgf("Request %v failed with %v", req, err)
			return nil, err
		}

		childSpan := opentracing.StartSpan(
			fmt.Sprintf("Size of %d bulk: %d", i, len(bulks[i])),
			opentracing.ChildOf(span.Context()),
		)
		childSpan.Finish()

		response.MeetingIds = append(response.MeetingIds, meetingIds...)
	}

	log.Printf("Сreation of the meetings was successful")

	return response, nil
}

func (a *api) CreateMeetingV1(
	ctx context.Context,
	req *desc.CreateMeetingV1Request,
) (*desc.CreateMeetingV1Response, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	meeting := &models.Meeting{
		UserId: req.Meeting.UserId,
		Link:   req.Meeting.Link,
		Start:  req.Meeting.Start.AsTime(),
		End:    req.Meeting.End.AsTime(),
	}

	err := a.repo.Add(ctx, meeting)
	if err != nil {
		log.Error().Msgf("Request %v failed with %v", req, err)
		return nil, err
	}

	log.Printf("Сreation of the meeting was successful")
	metrics.CreateCounterInc()

	event := producer.EventMessage{
		MeetingId: meeting.Id,
		Timestamp: time.Now().Unix(),
	}

	msg := producer.CreateMessage(producer.Create, event)

	go a.sendMessage(msg)

	return &desc.CreateMeetingV1Response{
		MeetingId: meeting.Id,
	}, nil
}

func (a *api) DescribeMeetingV1(
	ctx context.Context,
	req *desc.DescribeMeetingV1Request,
) (*desc.DescribeMeetingV1Response, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	meeting, err := a.repo.Describe(ctx, req.MeetingId)

	if err != nil {
		log.Error().Msgf("Request %v failed with %v", req, err)
		return nil, err
	}
	log.Printf("Reading of the meeting was successful")
	return &desc.DescribeMeetingV1Response{
		Meeting: &desc.Meeting{
			Id:     meeting.Id,
			UserId: meeting.UserId,
			Link:   meeting.Link,
			Start:  timestamppb.New(meeting.Start),
			End:    timestamppb.New(meeting.End),
		},
	}, nil
}

func (a *api) ListMeetingV1(
	ctx context.Context,
	req *desc.ListMeetingV1Request,
) (*desc.ListMeetingV1Response, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	meetings, err := a.repo.List(ctx, req.Limit, req.Offset)

	if err != nil {
		log.Error().Msgf("Request %v failed with %v", req, err)
		return nil, err
	}

	meetingList := make([]*desc.Meeting, 0, len(meetings))

	for _, meeting := range meetings {
		meetingList = append(meetingList, &desc.Meeting{
			Id:     meeting.Id,
			UserId: meeting.UserId,
			Link:   meeting.Link,
			Start:  timestamppb.New(meeting.Start),
			End:    timestamppb.New(meeting.End),
		})
	}
	log.Printf("Reading of the meetings was successful")
	return &desc.ListMeetingV1Response{
		Meetings: meetingList,
	}, nil
}

func (a *api) UpdateMeetingV1(
	ctx context.Context,
	req *desc.UpdateMeetingV1Request,
) (*desc.UpdateMeetingV1Response, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	meeting := models.Meeting{
		Id:     req.Meeting.Id,
		UserId: req.Meeting.UserId,
		Link:   req.Meeting.Link,
		Start:  req.Meeting.Start.AsTime(),
		End:    req.Meeting.End.AsTime(),
	}

	_, err := a.repo.Update(ctx, meeting)
	if err != nil {
		log.Error().Msgf("Request %v failed with %v", req, err)
		return nil, err
	}

	log.Printf("Updating of the meeting was successful")

	metrics.UpdateCounterInc()

	event := producer.EventMessage{
		MeetingId: meeting.Id,
		Timestamp: time.Now().Unix(),
	}
	msg := producer.CreateMessage(producer.Update, event)

	go a.sendMessage(msg)

	return &desc.UpdateMeetingV1Response{}, nil
}

func (a *api) RemoveMeetingV1(
	ctx context.Context,
	req *desc.RemoveMeetingV1Request,
) (*desc.RemoveMeetingV1Response, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	removed, err := a.repo.Remove(ctx, req.MeetingId)
	if err != nil {
		log.Error().Msgf("Request %v failed with %v", req, err)
		return nil, err
	}

	if removed == true {
		log.Printf("Removing of the meeting was successful")
	} else {
		log.Error().Msgf("Id was not found")
		return nil, errors.New("NotFound")
	}

	metrics.RemoveCounterInc()

	event := producer.EventMessage{
		MeetingId: req.MeetingId,
		Timestamp: time.Now().Unix(),
	}
	msg := producer.CreateMessage(producer.Delete, event)

	go a.sendMessage(msg)

	return &desc.RemoveMeetingV1Response{}, nil
}

func (a *api) sendMessage(msg producer.Message) {
	if err := a.prod.Send(msg); err != nil {
		log.Error().Msgf("failed send message to kafka")
	}
}
