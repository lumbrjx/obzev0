package main

import (
	"os"

	"github.com/sirupsen/logrus"
)

func Logger() *logrus.Entry {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
	logger.SetOutput(os.Stdout)

	rpcLogger = logger.WithFields(logrus.Fields{
		"service": "gRPC/server",
	})
	return rpcLogger

}
