package cron

import (
	"log"
	"time"

	"github.com/azinudinachzab/scr-syky-tech-test/service"
	"github.com/go-co-op/gocron"
)

type Cron struct {
	service service.Service
}

func NewCron(svc service.Service, loc *time.Location) *gocron.Scheduler {
	s := gocron.NewScheduler(loc)
	c := &Cron{
		service: svc,
	}

	s.WaitForScheduleAll()
	s.Every(5).Second().Do(c.FiveSecLogger)

	return s
}

func (c *Cron) FiveSecLogger() {
	log.Println("5 sec log")
}
