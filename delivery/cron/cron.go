package cron

import (
	"context"
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
	// s.Every(5).Second().Do(c.FiveSecLogger)
	s.Every(1).Day().At("00:00").Do(c.DailyBirthdayPromo)

	return s
}

//func (c *Cron) FiveSecLogger() {
//	log.Println("5 sec log")
//}

func (c *Cron) DailyBirthdayPromo() {
	ctx := context.Background()
	execTime := time.Now().Format("2006-01-02")
	log.Println("Executing job for date: ", execTime)
	c.service.SendDailyBirthdayPromo(ctx, time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(),
		0, 0, 0, 0, time.Now().Location()))
	log.Println("Job done for date: ", execTime)
}
