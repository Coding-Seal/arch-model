package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	journal_reader "github.com/Coding-Seal/arch-model/internal/journal_reader"
	"github.com/Coding-Seal/arch-model/internal/journal_reader/ui"
)

const fileName = ".journal/2024-12-02T01:40:25+03:00.jrl"

// const fileName = "../../.journal/2024-12-02T01:40:25+03:00.jrl"

const (
	numDoctors = 5
	benchCap   = 5
	numLobbies = 3
)

func main() {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatalln(err)
	}

	j := journal_reader.NewJournal(file)

	data, err := j.Read()
	if err != nil {
		log.Fatalln(err)
	}

	templ := template.Must(template.ParseFiles("cmd/journal_reader/state.html"))
	// templ := template.Must(template.ParseFiles("state.html"))

	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("GET /state", func(w http.ResponseWriter, r *http.Request) {
		stateID, err := strconv.Atoi(r.URL.Query().Get("id"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)

			return
		}

		if stateID >= len(data) || stateID < 0 {
			http.Error(w, "No such state", http.StatusNotFound)

			return
		}

		var eventStr string

		if stateID >= 0 {
			b, err := json.MarshalIndent(data[stateID], "", "    ")
			if err != nil {
				log.Fatalln(err)
			}

			eventStr = string(b)
			fmt.Println(eventStr)
		}

		system := ui.NewSystem(numDoctors, numLobbies, benchCap, stateID, eventStr)

		for _, e := range data[:stateID+1] {
			system.ApplyEvent(e)
		}

		err = templ.Execute(w, system)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
	http.ListenAndServe(":8080", nil)
}
