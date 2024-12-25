package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/Coding-Seal/arch-model/internal/bench"
	"github.com/Coding-Seal/arch-model/internal/doctor"
	"github.com/Coding-Seal/arch-model/internal/domain"
	event_manager "github.com/Coding-Seal/arch-model/internal/event_manager"
	"github.com/Coding-Seal/arch-model/internal/lobby"
	"github.com/Coding-Seal/arch-model/internal/nurse"
	"github.com/Coding-Seal/arch-model/pkg/jsonl"
)

const (
	firstAppointment = 20 * time.Millisecond
	genRate          = 20 * time.Millisecond
)

type Service interface {
	Run(ctx context.Context)
	Stop()
}

type simulationArgs struct {
	numDoctors int
	numSeats   int
	numLobbies int
	timeToRun  time.Duration
}

var runs = []simulationArgs{
	{numDoctors: 5, numSeats: 5, numLobbies: 3, timeToRun: 10 * time.Second},
	{numDoctors: 4, numSeats: 400, numLobbies: 3, timeToRun: 10 * time.Second},
	{numDoctors: 3, numSeats: 800, numLobbies: 3, timeToRun: 10 * time.Second},
	{numDoctors: 53, numSeats: 20, numLobbies: 33, timeToRun: 10 * time.Second},
	{numDoctors: 53, numSeats: 50, numLobbies: 33, timeToRun: 10 * time.Second},
}

func main() {
	// go func() { fmt.Println(http.ListenAndServe(":8080", nil)) }()
	wg := &sync.WaitGroup{}
	wg.Add(len(runs))

	for i, args := range runs {
		go func() {
			defer wg.Done()
			runSimulation(i, args)
			fmt.Println("Finished simulation", i, args)
		}()
	}

	wg.Wait()
	// time.Sleep(time.Hour)
}

func runSimulation(simulationNum int, args simulationArgs) {
	logFile, err := os.Create(".logs/" + strconv.Itoa(simulationNum) + ".log")
	if err != nil {
		log.Fatalln(err)
	}

	journalFile, err := os.Create(".journal/" + strconv.Itoa(simulationNum) + ".jrl")
	if err != nil {
		log.Fatalln(err)
	}

	parentCtx, _ := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	ctx, cancel := context.WithTimeout(parentCtx, args.timeToRun)
	defer cancel()

	log := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelDebug}))

	services := make([]Service, 0, args.numDoctors+args.numLobbies)
	bench := bench.New(log, args.numSeats)
	eventManager := event_manager.New(log, jsonl.NewWriter(journalFile))
	nurse := nurse.New(log, bench)
	nurse.Register(eventManager)

	seqID := domain.NewSeqID()
	for i := range args.numLobbies {
		lobby := lobby.New(log, bench, i+1, genRate, seqID)
		lobby.Register(eventManager)
		services = append(services, lobby)
	}

	for range args.numDoctors {
		doc := doctor.New(log, firstAppointment)
		doc.Register(eventManager, nurse)
		services = append(services, doc)
	}

	eventManager.Run(parentCtx, ctx)
	nurse.Run(parentCtx)

	for _, service := range services {
		service.Run(ctx)
	}

	<-ctx.Done()

	for _, service := range services {
		service.Stop()
	}

	nurse.Stop()
	eventManager.Stop()
}
