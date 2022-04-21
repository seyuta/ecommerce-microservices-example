package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"

	"github.com/seyuta/ecommerce-microservices-example/constant"
	"github.com/seyuta/ecommerce-microservices-example/s-transactions/repository"
	"github.com/seyuta/ecommerce-microservices-example/s-transactions/service"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func init() {
	// reading Command-Line arguments to determine configuration file
	args := os.Args[1:]
	var configname string = "s-transactions-config"
	if len(args) > 0 {
		configname = args[0] + "-config"
	}
	log.Printf("loading config file %s.yml", configname)
	viper.SetConfigName(configname)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fail loading config file: %v", err))
	}
	log.Println("config file loaded")

	// setup logging options
	logpath := viper.GetString(constant.CfgLogPath)
	log.Printf("setting log file path to %s", logpath)
	writer, err := rotatelogs.New(
		fmt.Sprintf("./log/%s.log", "%y%m%d"),
		rotatelogs.WithRotationTime(time.Hour*24),
	)
	if err != nil {
		panic(fmt.Errorf("failed to initialize log file: %v", err))
	}

	logrus.SetOutput(writer)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	logfmtr := &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			filename := path.Base(f.File)
			return "", fmt.Sprintf("%s %s:%d", time.Now().Format(time.RFC3339), filename, f.Line)
		},
	}
	logrus.SetFormatter(logfmtr)
}

func main() {
	logger := logrus.StandardLogger()

	initmongo := repository.InitMongoDB(logger)
	mongorepo := repository.BuildMongoRepository(initmongo)

	logger.Println("Starting up gRPC server...")
	server := grpc.NewServer(
		grpc.UnaryInterceptor(service.AuthInterceptor),
	)
	service.BuildGRPCService(server, mongorepo)

	port := viper.GetString(constant.CfgGrpcHost) + ":" + viper.GetString(constant.CfgGrpcPort)
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("could not listen to %s: %v", port, err)
	}

	logger.Printf("gRPC server on %s", port)

	panic(server.Serve(listener))
}
