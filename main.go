package main

import (
	"fmt"
	"os"

	"github.com/go-errors/errors"
	"github.com/goastro/indiapi/handler"
	"github.com/goastro/indiapi/routes"
	"github.com/goastro/indiclient"
	"github.com/goastro/indiserver"
	"github.com/rickbassham/goapi/middleware"
	"github.com/rickbassham/goapi/router"
	"github.com/rickbassham/goexec"
	"github.com/rickbassham/logging"
	"github.com/spf13/afero"
)

func main() {
	statusCode := mainWithStatusCode()
	os.Exit(statusCode)
}

func mainWithStatusCode() (statusCode int) {
	logger := logging.NewLogger(os.Stdout, logging.JSONFormatter{}, logging.LogLevelInfo)

	defer func() {
		if r := recover(); r != nil {
			statusCode = 1

			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}

			wrapped := errors.Wrap(err, 2)

			logger.WithField("stackTrace", wrapped.StackFrames()).WithError(err).Error("recovered from panic; exiting")
		}

		logger.WithField("statusCode", statusCode).Info("done")
	}()

	fs := afero.NewMemMapFs()

	requestLogger := middleware.NewRequestLogger(logger)

	svr := indiserver.NewINDIServer(logger, afero.NewOsFs(), "7624", goexec.ExecCommand{})

	svc := handler.NewINDIService(requestLogger, svr, indiclient.NetworkDialer{}, fs, 100)

	listenAddress := ":8080"

	indiRoutes := routes.NewINDIRoutes(svc)

	router := router.NewRouter(listenAddress, requestLogger, indiRoutes)

	logger.WithField("listenAddress", listenAddress).Info("listening...")

	err := router.ListenAndServe()
	if err != nil {
		logger.WithError(err).Error("error in router.ListenAndServe")
		statusCode = 1
		return
	}

	statusCode = 0
	return
}
