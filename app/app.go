package app

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/azinudinachzab/scr-syky-tech-test/delivery"
	"github.com/azinudinachzab/scr-syky-tech-test/model"
	"github.com/azinudinachzab/scr-syky-tech-test/pkg/clock"
	"github.com/azinudinachzab/scr-syky-tech-test/pkg/password"
	"github.com/azinudinachzab/scr-syky-tech-test/repository"
	"github.com/azinudinachzab/scr-syky-tech-test/service"
	"github.com/go-co-op/gocron"
	"github.com/go-playground/validator/v10"
)

type App struct {
	r             http.Handler
	s             *gocron.Scheduler
	conf          model.Configuration
	dbCore        *sql.DB
}

func New() *App {
	// init config
	conf := model.Configuration{
		AppAddress:         os.Getenv("APP_ADDRESS"),
		EmailUsername:      os.Getenv("OTP_ADDRESS"),
		EmailPassword:      os.Getenv("OTP_DURATION"),
		DatabaseName:       os.Getenv("DB_NAME"),
		CronTimezone:       os.Getenv("CRON_TIMEZONE"),
	}

	// set std log to print filename and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// init db
	db, err := dbConnection(conf.DatabaseName)
	if err != nil {
		log.Fatalln("can't connect database", err)
	}

	// init repo
	repo := repository.NewPgRepository(db)

	// validator
	v := validator.New()

	// clock
	clk := clock.NewTime()

	// password
	hash := password.New()

	// init service
	srv := service.NewService(service.Dependency{
		Validator: v,
		Repo:      repo,
		Clock:     clk,
		Hash:      hash,
		Conf:      conf,
	})
	// init http handler & cron
	timezone, err := time.LoadLocation(conf.CronTimezone)
	if err != nil {
		log.Fatalln("can't load timezone", err)
	}
	dlv := delivery.NewDelivery(delivery.Dependency{
		Service:  srv,
		Timezone: timezone,
	})

	return &App{r: dlv.HttpServer, s: dlv.Cron, conf: conf, dbCore: db}
}

func (a *App) Run() {
	// run server
	server := &http.Server{
		Addr:         a.conf.AppAddress,
		Handler:      a.r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Server run context
	go func() {
		log.Printf("server running on port %s", a.conf.AppAddress)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen and serve returned err: %v", err)
		}
	}()

	go func() {
		log.Printf("cron running")
		a.s.StartAsync()
	}()

	// Listen for syscall signals for process to interrupt/quit
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	// Shutdown signal with grace period of 15 seconds
	shutdownCtx, cancelShutdownCtx := context.WithTimeout(context.Background(), 15*time.Second)
	defer func() {
		cancelShutdownCtx()
		a.dbCore.Close()
		a.s.Stop()
	}()

	// Trigger graceful shutdown
	log.Println("server shutdown at: ", time.Now().Format(time.RFC3339))
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown err: %v", err)
	}
}
