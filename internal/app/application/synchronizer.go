package application

import (
	"person-key-saver/internal/app/builder"
	"person-key-saver/internal/app/container"
	"person-key-saver/internal/app/model"
	"person-key-saver/internal/app/repository"
	"sync"
)

const (
	GOROUTINES_LIMIT = 100
)

type Synchronizer struct {
	container *container.Container
}

func NewSynchronizer(container *container.Container) *Synchronizer {
	return &Synchronizer{
		container: container,
	}
}

func (this *Synchronizer) Run() {
	//logger := this.container.Get(container.LOGGER).(*logrus.Logger)
	keyRepository := this.container.Get(container.KEY_REPOSITORY).(*repository.PersonKeyRepository)
	deviceRepository := this.container.Get(container.DEVICE_REPOSITORY).(*repository.DeviceRepository)

	devices, _ := deviceRepository.GetItems()
	var wg sync.WaitGroup
	goroutines := make(chan int, GOROUTINES_LIMIT)
	for i := range devices {
		goroutines <- i
		wg.Add(1)
		go func(device *model.Device, goroutines <-chan int, wg *sync.WaitGroup) {
			syncBuilder := builder.NewSyncDataBuilder(this.container, device.Mac)
			data := syncBuilder.Build()
			for i := range data.Keys {
				keyRepository.Save(&data.Keys[i])
			}
			syncBuilder = nil
			<-goroutines
			wg.Done()
		}(devices[i], goroutines, &wg)
	}
	wg.Wait()
	close(goroutines)
}
