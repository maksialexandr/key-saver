package application

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"person-key-saver/internal/app/builder"
	"person-key-saver/internal/app/component"
	"person-key-saver/internal/app/config"
	"person-key-saver/internal/app/container"
	"person-key-saver/internal/app/form"
	"person-key-saver/internal/app/service"
)

const (
	EXCHANGE_NAME string = "*****"
	QUEUE_NAME    string = "*****"
	ROUTING_KEY   string = "*****"
)

type Listener struct {
	container *container.Container
	rmq       *component.Rmq
	config    config.RmqConfig
}

func NewListener(container *container.Container, rmq *component.Rmq, config config.RmqConfig) *Listener {
	return &Listener{
		container: container,
		rmq:       rmq,
		config:    config,
	}
}

func (this *Listener) Run() {
	reconnectEvent := make(chan bool, 1)
	logger := this.container.Get(container.LOGGER).(*logrus.Logger)
	go this.rmq.HandleReconnect(&this.config, EXCHANGE_NAME, QUEUE_NAME, ROUTING_KEY, reconnectEvent)
	for {
		select {
		case <-reconnectEvent:
			logger.Info("Registered handler after reconnect")
			go this.handleMessage(this.rmq.Messages)
		}
	}
}

// Обработчик сообщений amqp
func (this *Listener) handleMessage(messages <-chan amqp.Delivery) {
	logger := this.container.Get(container.LOGGER).(*logrus.Logger)
	keyService := this.container.Get(container.KEY_SERVICE).(*service.KeyService)
	for {
		select {
		case <-this.rmq.NotifyConnClose:
			return
		case message := <-messages:
			if len(message.Body) > 0 {
				go func() {
					payload := form.Payload{}
					if err := json.Unmarshal(message.Body, &payload); err != nil {
						logger.Error(err)
					}

					switch payload.Action {
					case container.ACTION_LOAD:
						dataBuilder := builder.NewSaveDataBuilder(this.container, &payload)
						data := dataBuilder.Build()
						for i := range data.Keys {
							_, err := keyService.Load(data.Device, data.Keys[i])
							if err != nil {
								logger.Error(err)
							}
						}
					case container.ACTION_ENCODE:
						dataBuilder := builder.NewEncodeDataBuilder(this.container, &payload)
						data := dataBuilder.Build()
						for i := range data.Keys {
							_, err := keyService.Encode(*data.Device, data.Keys[i])
							if err != nil {
								logger.Error(err)
							}
						}
					case container.ACTION_DECODE:
						dataBuilder := builder.NewDecodeDataBuilder(this.container, &payload)
						data := dataBuilder.Build()
						for i := range data.Keys {
							_, err := keyService.Decode(*data.Device, data.Keys[i])
							if err != nil {
								logger.Error(err)
							}
						}
					case container.ACTION_DELETE:
						dataBuilder := builder.NewDeleteDataBuilder(this.container, &payload)
						data := dataBuilder.Build()
						for i := range data.Keys {
							_, err := keyService.Delete(data.Device, data.Keys[i], payload.OnlyDb)
							if err != nil {
								logger.Error(err)
							}
						}
					}
					if err := message.Ack(false); err != nil {
						logger.Error(err)
					}
				}()
			}
		}
	}
}
