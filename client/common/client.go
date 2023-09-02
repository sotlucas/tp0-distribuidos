package common

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	ID            string
	ServerAddress string
	LoopLapse     time.Duration
	LoopPeriod    time.Duration
}

type UserBet struct {
	Nombre     string
	Apellido   string
	Documento  string
	Nacimiento string
	Numero     string
}

// Client Entity that encapsulates how
type Client struct {
	config ClientConfig
	conn   net.Conn
}

var signalChan chan (os.Signal) = make(chan os.Signal, 1)

// NewClient Initializes a new client receiving the configuration
// as a parameter
func NewClient(config ClientConfig) *Client {
	signal.Notify(signalChan, syscall.SIGTERM)
	client := &Client{
		config: config,
	}
	go client.shutdownClient()
	return client
}

// CreateClientSocket Initializes client socket. In case of
// failure, error is printed in stdout/stderr and exit 1
// is returned
func (c *Client) createClientSocket() error {
	conn, err := net.Dial("tcp", c.config.ServerAddress)
	if err != nil {
		log.Fatalf(
			"action: connect | result: fail | client_id: %v | error: %v",
			c.config.ID,
			err,
		)
	}
	c.conn = conn
	return nil
}

// Graceful shutdown of the client
func (c *Client) shutdownClient() {
	<-signalChan
	log.Debugf("action: shutdown_client | result: in_progress | client_id: %v", c.config.ID)
	c.conn.Close()
	log.Debugf("action: shutdown_client | result: success | client_id: %v", c.config.ID)
	os.Exit(0)
}

func (c *Client) sendUserBet(userBet UserBet) {
	// TODO: Modify the send to avoid short-write
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
	// Pad length to 4 bytes
	for len(lengthBytes) < 4 {
		lengthBytes = append([]byte("0"), lengthBytes...)
	}

	payloadBytes = append(lengthBytes, payloadBytes...)
	c.conn.Write(payloadBytes)
}

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop(userBet UserBet) {
	// autoincremental msgID to identify every message sent
	msgID := 1

loop:
	// Send messages if the loopLapse threshold has not been surpassed
	for timeout := time.After(c.config.LoopLapse); ; {
		select {
		case <-timeout:
			log.Infof("action: timeout_detected | result: success | client_id: %v",
				c.config.ID,
			)
			break loop
		default:
		}

		// Create the connection the server in every loop iteration. Send an
		c.createClientSocket()

		c.sendUserBet(userBet)
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		msgID++
		c.conn.Close()

		if err != nil {
			log.Errorf("action: receive_message | result: fail | client_id: %v | error: %v",
				c.config.ID,
				err,
			)
			return
		}
		log.Infof("action: receive_message | result: success | client_id: %v | msg: %v",
			c.config.ID,
			msg,
		)

		if msg == "OK\n" {
			log.Infof("action: apuesta_enviada | result: success | dni: %v | numero: %v",
				userBet.Documento,
				userBet.Numero,
			)
		} else {
			log.Errorf("action: apuesta_enviada | result: fail | dni: %v | numero: %v",
				userBet.Documento,
				userBet.Numero,
			)
		}

		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
