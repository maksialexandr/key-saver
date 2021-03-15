package component

import (
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"person-key-saver/internal/app/config"
	"sync"
	"time"
)

const (
	reconnectDelay        = 5 * time.Second
	reInitDelay           = 2 * time.Second
	PrefetchCountQos int  = 100
	PrefetchSizeQos  int  = 0
	GlobalQos        bool = false
)

var (
	errAlreadyClosed = errors.New("already closed: not connected to the server")
)

type Rmq struct {
	conn *amqp.Connection
	ch   *amqp.Channel

	eName    string
	qName    string
	rKey     string
	Messages <-chan amqp.Delivery

	mutex           sync.Mutex
	ready           bool
	NotifyConnClose chan *amqp.Error
	notifyChanClose chan *amqp.Error
}

func NewRmq() *Rmq {
	return &Rmq{}
}

func (r *Rmq) HandleReconnect(config *config.RmqConfig, eName string, qName string, rKey string, reconnectEvent chan bool) {
	r.eName = eName
	r.qName = qName
	r.rKey = rKey

	for {
		r.setReady(false)
		var err error
		r.conn, err = r.Connect(config)
		if err != nil {
			fmt.Printf("Failed to connect. Retrying...\n")
			<-time.After(reconnectDelay)
			continue
		}

		r.handleReInit(reconnectEvent)
	}
}

func (r *Rmq) handleReInit(reconnectEvent chan bool) bool {
	for {
		r.setReady(false)

		err := r.Init()

		r.Messages, err = r.Consume()
		if err != nil {
			return false
		}
		// отправляем событие что произошла успешная реконнекция
		reconnectEvent <- true

		if err != nil {
			<-time.After(reInitDelay)
			continue
		}

		select {
		case <-r.NotifyConnClose:
			fmt.Printf("Connection closed. Reconnecting...\n")
			return false
		case <-r.notifyChanClose:
			fmt.Printf("Channel closed. Re-running init...\n")
		}
	}
}

func (r *Rmq) Init() error {
	var err error

	r.ch, err = r.conn.Channel()

	if err != nil {
		return err
	}

	if err = r.ch.ExchangeDeclare(
		r.eName, // name
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	if _, err := r.ch.QueueDeclare(
		r.qName,
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	if err := r.ch.QueueBind(
		r.qName,
		r.rKey,
		r.eName,
		false,
		nil); err != nil {
		return err
	}

	if err := r.ch.Qos(PrefetchCountQos, PrefetchSizeQos, GlobalQos); err != nil {
		return err
	}

	r.changeChannel(r.ch)
	r.setReady(true)
	return nil
}

func (r *Rmq) changeChannel(channel *amqp.Channel) {
	r.mutex.Lock()
	r.ch = channel
	r.notifyChanClose = make(chan *amqp.Error)
	r.ch.NotifyClose(r.notifyChanClose)
	r.mutex.Unlock()
}

func (r *Rmq) Connect(config *config.RmqConfig) (*amqp.Connection, error) {
	conn, err := amqp.Dial(r.ConnectionFormat(config))

	if err != nil {
		return nil, err
	}

	r.changeConnection(conn)
	fmt.Printf("Connected AMQP!\n")
	return conn, nil
}

func (r *Rmq) changeConnection(connection *amqp.Connection) {
	r.mutex.Lock()
	r.conn = connection
	r.NotifyConnClose = make(chan *amqp.Error)
	r.conn.NotifyClose(r.NotifyConnClose)
	r.mutex.Unlock()
}

func (r *Rmq) setReady(state bool) {
	r.mutex.Lock()
	r.ready = state
	r.mutex.Unlock()
}

func (r *Rmq) isReady() bool {
	r.mutex.Lock()
	ready := r.ready
	r.mutex.Unlock()
	return ready
}

func (r *Rmq) Close() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !r.ready {
		return errAlreadyClosed
	}
	err := r.ch.Close()
	if err != nil {
		return err
	}
	err = r.conn.Close()
	if err != nil {
		return err
	}
	r.ready = false

	return nil
}

func (r *Rmq) Consume() (<-chan amqp.Delivery, error) {
	return r.ch.Consume(r.qName, "", false, false, false, false, nil)
}

func (r *Rmq) ConnectionFormat(config *config.RmqConfig) string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		config.User, config.Pwd, config.Host, config.Port, config.Vhost)
}
