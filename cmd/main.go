package main

import (
	"context"
	"medodstest/internal/server"
	"medodstest/pkg/handler"
	"medodstest/pkg/repository"
	"medodstest/pkg/service"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))

	dbconfig := repository.DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_DBNAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := repository.NewPostgresDB(dbconfig)
	if err != nil {
		logrus.Fatal(err)
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(server.Server)
	go func() {
		if err := srv.Run(os.Getenv("SERVER_PORT"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("Ошибка при запуске сервера: %s", err.Error())
		}
	}()

	logrus.Printf("Сервер прослушивает порт %s", os.Getenv("SERVER_PORT"))
	logrus.Print("Сервер запущен")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("Сервер выключается")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("Ошибка при выключении сервера: %s", err.Error())
	}

	db.Close()
}
