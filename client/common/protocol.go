package common

import (
	"encoding/binary"
	"fmt"
	"strings"
)

const LEN_BYTES int = 4
const BET_DELIMITER string = ";"

// Sends a slice of UserBet to the server
func (c *Client) sendUserBets(userBet []UserBet) {
	payloadBytes := serialize(c.config.ID, userBet)
	c.conn.Write(payloadBytes)
}

// Serializes a slice of UserBet into a byte array
func serialize(clientId string, userBet []UserBet) []byte {
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

	payloadBytes := []byte(payloadStr)
	payloadLength := len(payloadBytes)

	lengthBytes := make([]byte, LEN_BYTES)
	binary.BigEndian.PutUint32(lengthBytes, uint32(payloadLength))

	return append(lengthBytes, payloadBytes...)
}
