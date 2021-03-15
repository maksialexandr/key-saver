package application

import (
	"fmt"
	"*****/intercom/intercom-go-lib"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"log"
	"net/http"
	"person-key-saver/internal/app/application/configure"
	"person-key-saver/internal/app/application/route"
	"person-key-saver/internal/app/component"
	"person-key-saver/internal/app/container"
	"person-key-saver/internal/app/repository"
	"person-key-saver/internal/app/repository/api"
	"person-key-saver/internal/app/service"
	"time"
)

const (
	TIMER_EXECUTE time.Duration = time.Hour * 24
)

type Application struct {
	wsRouter   *route.WSRouter
	httpRouter *route.HttpRouter
	config     *configure.Config
	store      *component.Store
	container  *container.Container
}

func New() *Application {
	return &Application{
		wsRouter:   route.NewWSRouter(),
		httpRouter: route.NewHttpRouter(),
		store:      component.NewStore(),
		config:     configure.NewConfig(),
		container:  container.New(),
	}
}

func (this *Application) Configure(filename string) error {
	if err := this.config.SetUp(filename); err != nil {
		return err
	}

	logger, err := this.config.ConfigureLogger(this.config.LogLevel)
	if err != nil {
		return err
	}

	if err := this.store.Connect(this.config); err != nil {
		return err
	}
	logger.Info("Successfully connected to database!!!")

	var client intercom.ClientInterface
	if this.config.ClickHouse.IsValid() {
		err = intercom.OpenClickhouseConnection(&intercom.ClickhouseConfig{
			User:     this.config.ClickHouse.User,
			Pwd:      this.config.ClickHouse.Pwd,
			Dsn:      this.config.ClickHouse.Dsn,
			Port:     this.config.ClickHouse.Port,
			Database: this.config.ClickHouse.Database,
		})
		go intercom.LogWorker()

		if err != nil {
			return err
		}
		client = &intercom.LogClient{
			Client: &http.Client{},
		}
	} else {
		client = &intercom.SimpleClient{
			Client: &http.Client{},
		}
	}

	opts := mqtt.NewClientOptions().AddBroker(this.config.Mqtt.Host).SetClientID(this.config.Mqtt.Client)
	this.container.Set(container.LOGGER, logger)
	this.container.Set(container.WS_COMPONENT, component.NewWs())
	this.container.Set(container.DEVICE_API_REPOSITORY, api.NewDeviceRepository(client, mqtt.NewClient(opts)))
	this.container.Set(container.KEY_SERVICE, service.NewKeyService(this.container))
	this.container.Set(container.KEY_REPOSITORY, repository.NewPersonKeyRepository(this.store.Db))
	this.container.Set(container.DEVICE_REPOSITORY, repository.NewDeviceRepository(this.store.Db))
	return nil
}

func (this *Application) Run() {
	listener := NewListener(this.container, component.NewRmq(), this.config.Rmq)
	go listener.Run()

	//if this.config.DailyAutoSync {
	//	go this.startSynchronizer()
	//}

	this.wsRouter.Configure(this.container)
	go func() {
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", this.config.WS.Port), this.wsRouter.Router))
	}()

	this.httpRouter.Configure(this.container)
	go log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", this.config.Http.Port), this.httpRouter.Router))
}

func (this *Application) startSynchronizer() {
	synchronizer := NewSynchronizer(this.container)
	synchronizer.Run()
	for range time.NewTicker(TIMER_EXECUTE).C {
		synchronizer.Run()
	}
}
