package utils

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

// UserBet Represents a bet made by a user
type UserBet struct {
	Nombre     string
	Apellido   string
	Documento  string
	Nacimiento string
	Numero     string
}

// Reads a batch of bets from a csv file
func GetBets(betsFilepath string, batchSize int) []UserBet {
	f, err := os.Open(betsFilepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	bets := make([]UserBet, 0)
	for i := 0; i < batchSize; i++ {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		userBet := UserBet{
			Nombre:     record[0],
			Apellido:   record[1],
			Documento:  record[2],
			Nacimiento: record[3],
			Numero:     record[4],
		}

		bets = append(bets, userBet)
	}
	return bets
}
