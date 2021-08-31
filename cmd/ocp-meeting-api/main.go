package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/ozoncp/ocp-meeting-api/internal/config"
	"github.com/ozoncp/ocp-meeting-api/internal/metrics"
	"github.com/ozoncp/ocp-meeting-api/internal/producer"
	"github.com/ozoncp/ocp-meeting-api/internal/repo"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

func runGRPC(ctx context.Context, config *config.Config, prod producer.Producer) error {
	listen, err := net.Listen("tcp", config.Grpc.Address)
	if err != nil {
		log.Error().Err(err).Msg("GRPC: Listen")
		return err
	}

	db, err := repo.NewDB(config)
	if err != nil {
		log.Error().Err(err).Msg("db connect failed")
		return err
	}
	defer db.Close()

	s := grpc.NewServer()
	desc.RegisterOcpMeetingApiServer(s, api.NewOcpMeetingApi(repo.NewRepo(db), prod))
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

func runJSON(ctx context.Context, config *config.Config) error {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := desc.RegisterOcpMeetingApiHandlerFromEndpoint(ctx, mux, config.Grpc.Address, opts)
	if err != nil {
		log.Error().Err(err).Msg("JSON: Register API handler")
		return err
	}

	srv := &http.Server{Addr: config.Rest.Address, Handler: mux}
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

func runMetricsServer(config *config.Config) error {
	mux := http.NewServeMux()
	mux.Handle(config.Prometheus.Uri, promhttp.Handler())

	srv := &http.Server{
		Addr:    config.Prometheus.Port,
		Handler: mux,
	}
	metrics.RegisterMetrics()
	log.Info().Msg("Metrics server started")

	return srv.ListenAndServe()
}

func runKafka(config *config.Config) producer.Producer {
	prod, err := producer.NewProducer(config.Kafka.Topic, config.Kafka.Brokers)
	if err != nil {
		log.Panic().Msg("Kafka start failed")
	}
	log.Info().Msg("start producer")
	return prod
}

func main() {
	ctx := regSignalHandler(context.Background())

	config, err := config.Read()

	prod := runKafka(config)

	if err != nil {
		log.Fatal().Err(err).Msg("Readig configuration file was failed")
		return
	}

	go runMetricsServer(config)

	go func() {
		if err := runJSON(ctx, config); err != nil {
			log.Fatal().Err(err).Msg("HTTP Service stopped on error")
		}
	}()

	if err := runGRPC(ctx, config, prod); err != nil {
		log.Fatal().Err(err).Msg("GRPC Service stopped on error")
	}
}
