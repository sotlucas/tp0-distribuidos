package common

import (
	"bufio"
	"net"
	"os"
	"os/signal"
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
	// TODO: loopear todos los batches
	bets := utils.GetBets(c.config.BetsFilepath, c.config.BatchSize)

	c.createClientSocket()

	c.sendUserBets(bets)
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
		log.Infof("action: apuestas_enviadas | result: success")
	} else {
		log.Errorf("action: apuestas_enviadas | result: fail")
	}

	// Wait a time between sending one message and the next one
	time.Sleep(c.config.LoopPeriod)

	log.Infof("action: loop_finished | result: success | client_id: %v", c.config.ID)
}
