package common

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/utils"
)

const LEN_BYTES int = 4
const BET_DELIMITER string = ";"

// Sends a slice of Bet to the server
func (c *Client) sendBets(bet []utils.Bet) {
	payloadBytes := serialize(c.config.ID, bet)
	c.conn.Write(payloadBytes)
}

// Serializes a slice of Bet into a byte array
func serialize(clientId string, bets []utils.Bet) []byte {
	payloadBytes := buildPayload(bets, clientId)
	lengthBytes := buildLength(len(payloadBytes))
	return append(lengthBytes, payloadBytes...)
}

// Builds the payload of the message into a byte array
func buildPayload(bets []utils.Bet, clientId string) []byte {
	payload := make([]string, 0, len(bets))
	for _, bet := range bets {
		betStr := fmt.Sprintf(
			"%s:%s:%s:%s:%s:%s",
			clientId,
			bet.Nombre,
			bet.Apellido,
			bet.Documento,
			bet.Nacimiento,
			bet.Numero,
		)
		payload = append(payload, betStr)
	}
	payloadStr := strings.Join(payload, BET_DELIMITER)
	return []byte(payloadStr)
}

// Builds the length of the payload into a byte array
func buildLength(payloadLength int) []byte {
	lengthBytes := make([]byte, LEN_BYTES)
	binary.BigEndian.PutUint32(lengthBytes, uint32(payloadLength))
	return lengthBytes
}
