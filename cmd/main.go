package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/vpnvsk/p_s"
	handl "github.com/vpnvsk/p_s/pkg/handler"
	"github.com/vpnvsk/p_s/pkg/repository"
	"github.com/vpnvsk/p_s/pkg/service"
)

// @title Credentials Manager API
// @version 1.0
// @description API Server for Credentials Manager Application

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	if err := initConfig(); err != nil {
		logrus.Fatalf("error while reading config files %s", err.Error())
	}
	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("Failed to load config: %s", err.Error())
	}
	db, err := repository.NewPostgresDb(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	})
	if err != nil {
		logrus.Fatalf("Failed to connect to database:%s", err.Error())
	}
	serviceConfig := service.Config{}
	strAppId, err := strconv.ParseInt(os.Getenv("appId"), 10, 32)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	appId := int32(strAppId)
	serviceConfig.SetFields(
		os.Getenv("salt"),
		os.Getenv("signingKey"),
		os.Getenv("key"),
		2*time.Hour,
		appId,
	)
	repos := repository.NewRepository(db)
	services := service.NewService(repos, serviceConfig)
	handler := handl.NewHandler(services)
	srv := new(p_s.Server)
	go func() {
		if err := srv.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
			logrus.Fatalf("error while running server %s", err.Error())
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit
	if err = srv.ShutDown(context.Background()); err != nil {
		logrus.Errorf("error while shuting down: %s", err.Error())
	}
	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
