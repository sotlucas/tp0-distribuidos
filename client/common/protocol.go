package common

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/utils"
	log "github.com/sirupsen/logrus"
)

const LEN_BYTES int = 4
const BET_DELIMITER string = ";"

// Sends a slice of Bet to the server
func (c *Client) sendBets(bet []utils.Bet) {
	c.createClientSocket()

	payloadBytes := serialize(c.config.ID, bet)
	sz, err := c.conn.Write(payloadBytes)
	if err != nil {
		log.Fatalf(
			"action: send | result: fail | client_id: %v | sz: %v | error: %v",
			c.config.ID,
			sz,
			err,
		)
	}
	log.Debugf(
		"action: send | result: success | client_id: %v | sz: %v | payload_size: %v | payload: %v",
		c.config.ID,
		sz,
		len(payloadBytes),
		string(payloadBytes),
	)

	msg, err := bufio.NewReader(c.conn).ReadString('\n')
	c.conn.Close()

	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}
	log.Debugf("action: receive_message | result: success | client_id: %v | msg: %v",
		c.config.ID,
		msg,
	)

	if msg == "OK\n" {
		log.Infof("action: apuestas_enviadas | result: success")
	} else {
		log.Errorf("action: apuestas_enviadas | result: fail")
	}
}

// Sends a FINISH message to the server indicating that the client
// has finished sending bets
func (c *Client) sendFinish() {
	c.createClientSocket()

	action := []byte("FINISH::")
	lengthBytes := buildLength(len(action))
	msg := append(lengthBytes, action...)
	sz, err := c.conn.Write(msg)
	if err != nil {
		log.Fatalf(
			"action: send | result: fail | client_id: %v | sz: %v | error: %v",
			c.config.ID,
			sz,
			err,
		)
	}

	res, err := bufio.NewReader(c.conn).ReadString('\n')
	c.conn.Close()

	if err != nil {
		log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
		return
	}
	log.Debugf("action: receive_message | result: success | client_id: %v | msg: %v",
		c.config.ID,
		msg,
	)

	if res == "OK\n" {
		log.Infof("action: fin_envio_apuestas | result: success")
	} else {
		log.Errorf("action: fin_envio_apuestas | result: fail")
	}
}

// Serializes a slice of Bet into a byte array
func serialize(clientId string, bets []utils.Bet) []byte {
	action := []byte("BET::")
	payloadBytes := buildPayload(bets, clientId)
	message := append(action, payloadBytes...)
	lengthBytes := buildLength(len(message))
	return append(lengthBytes, message...)
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
