package api

import (
	"context"
	desc "github.com/ozoncp/ocp-meeting-api/pkg/ocp-meeting-api"
	log "github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type api struct {
	desc.UnimplementedOcpMeetingApiServer
}

func NewOcpMeetingApi() desc.OcpMeetingApiServer {
	return &api{}
}

func (a *api) CreateMeetingV1(
	ctx context.Context,
	req *desc.CreateMeetingV1Request,
) (*desc.CreateMeetingV1Response, error) {
	log.Printf("Сreation of the meeting was successful")
	return &desc.CreateMeetingV1Response{}, nil
}

func (a *api) DescribeMeetingV1(
	ctx context.Context,
	req *desc.DescribeMeetingV1Request,
) (*desc.DescribeMeetingV1Response, error) {
	log.Printf("Reading of the meeting was successful")
	return &desc.DescribeMeetingV1Response{}, nil
}

func (a *api) ListMeetingV1(
	ctx context.Context,
	req *emptypb.Empty,
) (*desc.ListMeetingV1Response, error) {
	log.Printf("Reading of the meetings was successful")
	return &desc.ListMeetingV1Response{}, nil
}

func (a *api) RemoveMeetingV1(
	ctx context.Context,
	req *desc.RemoveMeetingV1Request,
) (*desc.RemoveMeetingV1Response, error) {
	log.Printf("Removing of the meeting was successful")
	return &desc.RemoveMeetingV1Response{}, nil
}
