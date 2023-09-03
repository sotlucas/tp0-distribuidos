package common

import (
	"bufio"
	"encoding/csv"
	"io"
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
	LoopPeriod    time.Duration
	BetsFilepath  string
	BatchSize     int
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

// StartClientLoop Send messages to the client until some time threshold is met
func (c *Client) StartClientLoop() {
	for _, userBet := range getBets(c.config.BetsFilepath, c.config.BatchSize) {
		c.createClientSocket()

		c.sendUserBet(userBet)
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
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

// Reads a batch of bets from a csv file
func getBets(betsFilepath string, batchSize int) []UserBet {
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
