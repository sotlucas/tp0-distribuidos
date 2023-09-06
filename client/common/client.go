package common

import (
	"encoding/csv"
	"io"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/7574-sistemas-distribuidos/docker-compose-init/client/utils"
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
	f, err := os.Open(c.config.BetsFilepath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)

	hasNext := true
	for hasNext {
		bets, next := nextBatch(csvReader, c.config.BatchSize)
		hasNext = next

		c.sendBets(bets)

		// Wait a time between sending one message and the next one
		time.Sleep(c.config.LoopPeriod)
	}

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
	c.sendFinish()

	for {
		msg := c.askWinner()

		if msg.Action == "WINNER" {
			cantGanadores := len(strings.Split(msg.Payload, ";"))
			log.Infof("action: consulta_ganadores | result: success | client_id: %v | cant_ganadores: %v", c.config.ID, cantGanadores)
			break
		} else if msg.Action == "WINNERWAIT" {
			log.Infof("action: winner | result: in_progress | client_id: %v | msg: %v", c.config.ID, msg.Payload)
			waitTime, _ := time.ParseDuration(msg.Payload + "s")
			time.Sleep(waitTime) // TODO: modificar para que no sea busy wait
		} else {
			log.Errorf("action: winner | result: fail | client_id: %v", c.config.ID)
			break
		}

	}
}

func nextBatch(csvReader *csv.Reader, batchSize int) ([]utils.Bet, bool) {
	next := true
	bets := make([]utils.Bet, 0)
	for i := 0; i < batchSize; i++ {
		record, err := csvReader.Read()
		if err == io.EOF {
			next = false
			break
		}
		if err != nil {
			next = false
			log.Fatal(err)
		}

		bet := utils.Bet{
			Nombre:     record[0],
			Apellido:   record[1],
			Documento:  record[2],
			Nacimiento: record[3],
			Numero:     record[4],
		}

		bets = append(bets, bet)
	}
	return bets, next
}
