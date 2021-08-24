package api

import (
	"context"
	"github.com/ozoncp/ocp-meeting-api/internal/models"
	"github.com/ozoncp/ocp-meeting-api/internal/repo"
	desc "github.com/ozoncp/ocp-meeting-api/pkg/ocp-meeting-api"
	log "github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type api struct {
	desc.UnimplementedOcpMeetingApiServer
	repo repo.Repo
}

func NewOcpMeetingApi(repo repo.Repo) desc.OcpMeetingApiServer {
	return &api{
		repo: repo,
	}
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

	updated, err := a.repo.Update(ctx, meeting)
	if err != nil {
		log.Error().Msgf("Request %v failed with %v", req, err)
		return nil, err
	}

	log.Printf("Updating of the meeting was successful")
	return &desc.UpdateMeetingV1Response{
		Updated: updated,
	}, nil
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
		log.Printf("Removing of the meeting was failed")
	}

	return &desc.RemoveMeetingV1Response{
		Removed: removed,
	}, nil
}
