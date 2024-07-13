package interceptors

import (
	"runtime/debug"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sirupsen/logrus"
)

func RecoveryHandler(rpcLogger *logrus.Entry) func(p any) error {
	return func(p any) error {
		stack := debug.Stack()
		rpcLogger.WithFields(logrus.Fields{
			"panic": p,
			"stack": string(stack),
		}).Error("recovered from panic")
		return status.Errorf(codes.Internal, "%s", p)
	}
}
