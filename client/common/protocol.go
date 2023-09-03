package common

import (
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/utils"
)

const LEN_BYTES int = 4
const BET_DELIMITER string = ";"

// Sends a slice of UserBet to the server
func (c *Client) sendUserBets(userBet []utils.UserBet) {
	payloadBytes := serialize(c.config.ID, userBet)
	c.conn.Write(payloadBytes)
}

// Serializes a slice of UserBet into a byte array
func serialize(clientId string, userBet []utils.UserBet) []byte {
	payloadBytes := buildPayload(userBet, clientId)
	lengthBytes := buildLength(len(payloadBytes))
	return append(lengthBytes, payloadBytes...)
}

// Builds the payload of the message into a byte array
func buildPayload(userBet []utils.UserBet, clientId string) []byte {
	payload := make([]string, 0, len(userBet))
	for _, bet := range userBet {
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
