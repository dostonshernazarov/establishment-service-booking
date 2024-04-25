package main

import (
	"Booking/establishment-service-booking/internal/app"
	"Booking/establishment-service-booking/internal/pkg/config"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	// initialization config
	config := config.New()

	// initialization app
	app, err := app.NewApp(config)
	if err != nil {
		log.Fatal(err)
	}

	// Run the application
	go func() {
		if err := app.Run(); err != nil {
			app.Logger.Error("app run", zap.Error(err))
		}
	}()

	// graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	app.Logger.Info("User service stops !")

	// app stops
	app.Stop()

}
