package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Coding-Seal/arch-model/internal/bench"
	"github.com/Coding-Seal/arch-model/internal/doctor"
	event_manager "github.com/Coding-Seal/arch-model/internal/event_manager"
	"github.com/Coding-Seal/arch-model/internal/lobby"
	"github.com/Coding-Seal/arch-model/internal/nurse"
	"github.com/Coding-Seal/arch-model/pkg/jsonl"
)

const (
	numDoctors       = 5
	numberOfSeats    = 5
	numberOfLobbies  = 3
	firstAppointment = 1 * time.Second
	genRate          = 5 * time.Second
)

type Service interface {
	Run(ctx context.Context)
	Stop()
}

func main() {
	logFile, err := os.Create(".logs/" + time.Now().Format(time.RFC3339) + ".log")
	if err != nil {
		log.Fatalln(err)
	}

	journalFile, err := os.Create(".journal/" + time.Now().Format(time.RFC3339) + ".jrl")
	if err != nil {
		log.Fatalln(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	log := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug}))

	services := make([]Service, 0, numDoctors+numberOfLobbies)
	bench := bench.New(log, numberOfSeats)
	eventManager := event_manager.New(log, jsonl.NewWriter(journalFile))
	nurse := nurse.New(log, bench)
	nurse.Register(eventManager)

	services = append(services, eventManager, nurse)

	for i := range numberOfLobbies {
		lobby := lobby.New(log, bench, i+1, genRate)
		lobby.Register(eventManager)
		services = append(services, lobby)
	}

	for range numDoctors {
		doc := doctor.New(log, firstAppointment)
		doc.Register(eventManager, nurse)
		services = append(services, doc)
	}

	for _, service := range services {
		service.Run(ctx)
	}

	<-ctx.Done()

	for _, service := range services {
		service.Stop()
	}
}
