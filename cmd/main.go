package main

import (
	"context"
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
)

const (
	numDoctors       = 5
	numberOfSeats    = 5
	numberOfLobbies  = 3
	firstAppointment = time.Second
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	bench := bench.New(log, numberOfSeats)
	eventManager := event_manager.New(log)
	nurse := nurse.New(log, bench)

	nurse.Register(eventManager)

	lobbies := make([]*lobby.Lobby, 0, numberOfLobbies)

	doctors := make([]*doctor.Doctor, 0, numDoctors)

	for i := range numberOfLobbies {
		lobby := lobby.New(log, bench, i+1)
		lobby.Register(eventManager)
		lobbies = append(lobbies, lobby)
	}

	for range numDoctors {
		doc := doctor.New(firstAppointment)
		doc.Register(eventManager, nurse)
		doctors = append(doctors, doc)
	}

	eventManager.Run(ctx)
	nurse.Run(ctx)

	for _, doc := range doctors {
		doc.Run(ctx)
	}

	for _, lobby := range lobbies {
		lobby.Run(ctx)
	}

	<-ctx.Done()

}
