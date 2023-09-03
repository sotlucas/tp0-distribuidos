package utils

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

// Bet Represents a bet made by a user
type Bet struct {
	Nombre     string
	Apellido   string
	Documento  string
	Nacimiento string
	Numero     string
}

// Reads a batch of bets from a csv file
func GetBets(betsFilepath string, batchSize int) []Bet {
	f, err := os.Open(betsFilepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	bets := make([]Bet, 0)
	for i := 0; i < batchSize; i++ {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		bet := Bet{
			Nombre:     record[0],
			Apellido:   record[1],
			Documento:  record[2],
			Nacimiento: record[3],
			Numero:     record[4],
		}

		bets = append(bets, bet)
	}
	return bets
}
