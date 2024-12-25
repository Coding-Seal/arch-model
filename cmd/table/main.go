package main

import (
	"fmt"
	"io/fs"
	"log"
	"maps"
	"math"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/Coding-Seal/arch-model/internal/domain"
	journal_reader "github.com/Coding-Seal/arch-model/internal/journal_reader"
)

const dirName = ".journal"

func main() {
	err := filepath.WalkDir(dirName, printTable)
	if err != nil {
		log.Fatalln(err)
	}
}

func printTable(path string, info fs.DirEntry, err error) error {
	if err != nil {
		return err
	}

	if info.IsDir() {
		return nil
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	j := journal_reader.NewJournal(file)

	data, err := j.Read()
	if err != nil {
		log.Fatalln(err)
	}

	ticketMap := make(map[int]*Ticket, 0)

	for _, e := range data {
		// log.Printf("handel event: %+v\n", e)
		switch event := e.(type) {
		case domain.NewPatientEvent:
			ticketMap[event.Patient.ID] = &Ticket{
				patientID:     event.Patient.ID,
				lobbyID:       event.LobbyID,
				timeGenerated: event.Timestamp,
			}
		case domain.PatientInQueueEvent:
			t := ticketMap[event.PatientID]
			t.timeStartInQueue = event.Timestamp
		case domain.PatientGoneEvent:
			t := ticketMap[event.PatientID]
			t.timeEndInQueue = event.Timestamp
			t.finished = false
			t.denied = true
		case domain.AppointmentFinishedEvent:
			t := ticketMap[event.PatientID]
			t.timeEndWork = event.Timestamp
			t.finished = true
		case domain.AppointmentStartedEvent:
			t := ticketMap[event.PatientID]
			t.timeEndInQueue = event.Timestamp
			t.timeStartWork = event.Timestamp
			t.doctorID = event.DoctorID
		}
	}

	doctors := make(map[int]*Doctor)
	lobbies := make(map[int]*Lobby)

	for _, ticket := range ticketMap {
		doc, ok := doctors[ticket.doctorID]
		if !ok {
			doc = &Doctor{ID: ticket.doctorID}
			doctors[ticket.doctorID] = doc
		}

		lobby, ok := lobbies[ticket.lobbyID]
		if !ok {
			lobby = &Lobby{ID: ticket.lobbyID}
			lobbies[ticket.lobbyID] = lobby
		}

		doc.addTicket(ticket)
		lobby.addTicket(ticket)
	}

	fmt.Println(path)
	fmt.Println()
	fmt.Println("NUM_TICKETS:", len(ticketMap))
	fmt.Println("DOCTORS")

	docs := slices.SortedFunc(maps.Values(doctors), func(l, r *Doctor) int {
		return l.ID - r.ID
	})

	lob := slices.SortedFunc(maps.Values(lobbies), func(l, r *Lobby) int {
		return l.ID - r.ID
	})

	for _, doc := range docs {
		start := doc.timeFirstPatientStarted
		end := doc.timeLastPatientFinished
		timeBusy := doc.timeBusy

		ratio := math.Round(float64(timeBusy) / float64(end.Sub(start)) * 100)
		if math.IsNaN(ratio) {
			ratio = 0
		}

		fmt.Println("ID :", doc.ID, "busyRatio :", ratio, "%")
	}

	fmt.Println("LOBBIES")

	for _, lobby := range lob {
		ratio := math.Round(float64(lobby.numDenials) / float64(lobby.numGenerated) * 100)

		fmt.Println()

		inQueueMean := float64(lobby.timeInQueue.Milliseconds()) / float64(lobby.numGenerated)
		inQueueDisp := lobby.timeInQueueSQR/float64(lobby.numGenerated) - inQueueMean*inQueueMean
		inWorkMean := float64(lobby.timeInWork.Milliseconds()) / float64(lobby.numGenerated)
		inWorkDisp := lobby.timeInWorkSQR/float64(lobby.numGenerated) - inWorkMean*inWorkMean

		fmt.Println("ID :", lobby.ID, "denial/generated :", ratio, "%")
		fmt.Println("	in queue	", "mean: ", time.Duration(inQueueMean)*time.Millisecond, "disp: ", time.Duration(inQueueDisp)*time.Millisecond)
		fmt.Println("	in work		", "mean: ", time.Duration(inWorkMean)*time.Millisecond, "disp: ", time.Duration(inWorkDisp)*time.Millisecond)
	}

	fmt.Println()
	var cost int64
	for _, lobby := range lob {
		cost += int64(lobby.numDenials * 500)
		cost += int64((lobby.numGenerated - lobby.numDenials) * 250)

		cost += int64((lobby.timeInQueue + lobby.timeInWork).Milliseconds() * 20)
	}
	cost += int64(len(docs) * 30000)

	fmt.Println("COSTS: ", cost)
	fmt.Println()

	return nil
}

type Ticket struct {
	patientID int
	lobbyID   int
	doctorID  int

	denied           bool
	finished         bool
	timeGenerated    time.Time
	timeStartInQueue time.Time
	timeEndInQueue   time.Time
	timeStartWork    time.Time
	timeEndWork      time.Time
}

type Doctor struct {
	ID int

	timeFirstPatientStarted time.Time
	timeLastPatientFinished time.Time
	timeBusy                time.Duration
}

func (d *Doctor) addTicket(t *Ticket) {
	if !t.finished {
		return
	}

	if d.timeFirstPatientStarted.IsZero() {
		d.timeFirstPatientStarted = t.timeStartWork
		d.timeLastPatientFinished = t.timeEndWork
	} else if d.timeFirstPatientStarted.After(t.timeStartWork) {
		d.timeFirstPatientStarted = t.timeStartWork
	} else if d.timeLastPatientFinished.Before(t.timeEndWork) {
		d.timeLastPatientFinished = t.timeEndWork
	}

	d.timeBusy += t.timeEndWork.Sub(t.timeStartWork)
}

type Lobby struct {
	ID int

	numGenerated   int
	numDenials     int
	timeInQueue    time.Duration
	timeInQueueSQR float64
	timeInWork     time.Duration
	timeInWorkSQR  float64
}

func (l *Lobby) addTicket(t *Ticket) {
	l.numGenerated++

	if t.finished {
		delta := t.timeEndWork.Sub(t.timeStartWork)
		l.timeInWork += delta
		l.timeInWorkSQR += float64(delta.Milliseconds()) * float64(delta.Milliseconds())
	}
	if t.denied {
		l.numDenials++
	}

	if !t.timeEndInQueue.IsZero() {
		delta := t.timeEndInQueue.Sub(t.timeStartInQueue)
		l.timeInQueue += delta
		l.timeInQueueSQR += float64(delta.Milliseconds()) * float64(delta.Milliseconds())
	}
}
