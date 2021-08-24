package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/jmoiron/sqlx"
	"github.com/ozoncp/ocp-meeting-api/internal/db"
	"github.com/ozoncp/ocp-meeting-api/internal/repo"
	"google.golang.org/grpc"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	api "github.com/ozoncp/ocp-meeting-api/internal/api"
	desc "github.com/ozoncp/ocp-meeting-api/pkg/ocp-meeting-api"
	log "github.com/rs/zerolog/log"
)

const (
	grpcPort = ":8082"
	httpPort = ":8080"
)

func regSignalHandler(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer signal.Stop(done)
		<-done
		log.Info().Msg("Stop signal received")
		cancel()
	}()

	return ctx
}

func runGRPC(ctx context.Context, database *sqlx.DB) error {
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Error().Err(err).Msg("GRPC: Listen")
		return err
	}

	s := grpc.NewServer()
	desc.RegisterOcpMeetingApiServer(s, api.NewOcpMeetingApi(repo.NewRepo(database)))
	log.Info().Msg("GRPC Service was started")

	srvErr := make(chan error)
	go func() {
		if err := s.Serve(listen); err != nil {
			srvErr <- err
		}
	}()

	select {
	case err := <-srvErr:
		log.Error().Err(err).Msg("GRPC: Serve")
		return err

	case <-ctx.Done():
		s.GracefulStop()
		log.Info().Msg("GRPC Service was closed")
	}

	return nil
}

func runJSON(ctx context.Context) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := desc.RegisterOcpMeetingApiHandlerFromEndpoint(ctx, mux, grpcPort, opts)
	if err != nil {
		log.Error().Err(err).Msg("JSON: Register API handler")
		return err
	}

	srv := &http.Server{Addr: httpPort, Handler: mux}
	log.Info().Msg("HTTP Service was started")

	srvErr := make(chan error)
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			srvErr <- err
		}
	}()

	select {
	case err := <-srvErr:
		log.Error().Err(err).Msg("JSON: Serve")
		return err

	case <-ctx.Done():
		_ = srv.Shutdown(ctx)
		log.Info().Msg("HTTP Service was closed")
	}

	return nil
}

func main() {
	ctx := regSignalHandler(context.Background())

	database := db.Connect("postgres://postgres:postgres@127.0.0.1:5432/postgres?sslmode=disable")
	defer database.Close()

	go func() {
		if err := runJSON(ctx); err != nil {
			log.Fatal().Err(err).Msg("HTTP Service stopped on error")
		}
	}()

	if err := runGRPC(ctx, database); err != nil {
		log.Fatal().Err(err).Msg("GRPC Service stopped on error")
	}
}
