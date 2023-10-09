package delivery

import (
	"github.com/azinudinachzab/scr-syky-tech-test/delivery/cron"
	"net/http"
	"time"

	httpDelivery "github.com/azinudinachzab/scr-syky-tech-test/delivery/http"
	"github.com/azinudinachzab/scr-syky-tech-test/service"
	"github.com/go-co-op/gocron"
)

type(
	Dependency struct {
		Service  service.Service
		Timezone *time.Location
	}

	Delivery struct {
		HttpServer http.Handler
		Cron       *gocron.Scheduler
	}
)

func NewDelivery(dep Dependency) *Delivery {
	return &Delivery{
		HttpServer: httpDelivery.NewHttpServer(dep),
		Cron:       cron.NewCron(dep.Service, dep.Timezone),
	}
}
