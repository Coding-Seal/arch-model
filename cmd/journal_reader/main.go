package main

import (
	"fmt"
	"log"
	"os"
	"slices"

	"github.com/Coding-Seal/arch-model/internal/domain"
	journal_reader "github.com/Coding-Seal/arch-model/internal/journal_reader"
)

const fileName = ".journal/2024-11-28T23:57:03+03:00.jrl"

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

	slices.SortFunc(data, func(lhs, rhs domain.Event) int {
		return int(lhs.Time().Sub(rhs.Time()))
	})
	fmt.Println(data)
}
