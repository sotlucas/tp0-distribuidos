package common

import "fmt"

// Sends a slice of UserBet to the server
func (c *Client) sendUserBets(userBet []UserBet) {
	payloadBytes := serialize(c.config.ID, userBet)
	c.conn.Write(payloadBytes)
}

// Serializes a slice of UserBet into a byte array
func serialize(clientId string, userBet []UserBet) []byte {
	payload := ""
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
		payload += betStr + ";"
	}

	payloadBytes := []byte(payload)
	payloadLength := len(payloadBytes)

	lengthBytes := []byte(fmt.Sprintf("%d", payloadLength))
	// Pad length to 4 bytes with leading zeros
	for len(lengthBytes) < 4 {
		lengthBytes = append([]byte("0"), lengthBytes...)
	}

	return append(lengthBytes, payloadBytes...)
}
