package common

import "fmt"

func (c *Client) sendUserBet(userBet UserBet) {
	payload := fmt.Sprintf(
		"%s:%s:%s:%s:%s:%s",
		c.config.ID,
		userBet.Nombre,
		userBet.Apellido,
		userBet.Documento,
		userBet.Nacimiento,
		userBet.Numero,
	)

	payloadBytes := []byte(payload)
	payloadLength := len(payloadBytes)

	lengthBytes := []byte(fmt.Sprintf("%d", payloadLength))
	// Pad length to 4 bytes with leading zeros
	for len(lengthBytes) < 4 {
		lengthBytes = append([]byte("0"), lengthBytes...)
	}

	payloadBytes = append(lengthBytes, payloadBytes...)
	c.conn.Write(payloadBytes)
}
